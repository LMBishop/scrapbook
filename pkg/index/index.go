package index

import (
	"maps"
	"slices"
	"sort"
	"sync"

	"github.com/LMBishop/scrapbook/pkg/site"
)

type SiteIndex struct {
	mu          sync.RWMutex
	sites       map[string]*site.Site
	sitesByHost map[string]*site.Site
}

func NewSiteIndex() *SiteIndex {
	var siteIndex SiteIndex
	siteIndex.sites = make(map[string]*site.Site)
	siteIndex.sitesByHost = make(map[string]*site.Site)
	return &siteIndex
}

func (s *SiteIndex) GetSiteByHost(host string) *site.Site {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sitesByHost[host]
}

func (s *SiteIndex) GetSite(site string) *site.Site {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sites[site]
}

func (s *SiteIndex) GetSites() []*site.Site {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sites := slices.Collect(maps.Values(s.sites))
	sort.Slice(sites, func(i, j int) bool {
		return sites[i].Name < sites[j].Name
	})
	return sites
}

func (s *SiteIndex) AddSite(site *site.Site) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sites[site.Name] = site
	s.updateSiteIndexes()
}

func (s *SiteIndex) RemoveSite(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sites, name)
	s.updateSiteIndexes()
}

func (s *SiteIndex) UpdateSiteIndexes() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.updateSiteIndexes()
}

func (s *SiteIndex) updateSiteIndexes() {
	clear(s.sitesByHost)
	for _, site := range s.sites {
		if site.SiteConfig.Host != "" {
			s.sitesByHost[site.SiteConfig.Host] = site
		}
	}
}
