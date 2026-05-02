package html

import (
	"strconv"
	"time"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/site"
	. "github.com/LMBishop/scrapbook/web/skeleton"
	"github.com/dustin/go-humanize"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func VersionPage(site *site.Site, versionMeta config.VersionMeta) Node {
	createDate := time.Unix(int64(versionMeta.Created), 0)
	currentVersion, _ := site.GetCurrentVersion()

	return Page("Version "+versionMeta.Hash[0:8],
		H1(Text("Version "+versionMeta.Hash[0:8])),

		FieldSet(
			Legend(Text("Version actions")),

			Div(
				Class("control-group"),

				If(currentVersion != versionMeta.Hash,
					Form(
						Method("POST"),
						Action("current"),

						Input(
							Type("submit"),
							Value("Set as current version"),
						),
					),
				),
				NavButton("Download version", "#"),
				NavButton("Delete version", "delete"),
			),
		),

		P(
			Dl(
				Dt(Text("Hash")),
				Dd(Text(versionMeta.Hash)),

				Dt(Text("Upload kind")),
				Dd(Text(versionMeta.Kind)),

				Dt(Text("Original name")),
				Dd(Text(versionMeta.Original)),

				Dt(Text("Number of files")),
				Dd(Text(strconv.Itoa(int(versionMeta.Files)))),

				Dt(Text("Size on disk")),
				Dd(Text(humanize.IBytes(uint64(versionMeta.Size)))),

				Dt(Text("Created at")),
				Dd(Text(createDate.Format(time.UnixDate))),

				Dt(Text("Source")),
				Dd(Text(versionMeta.Source)),

				Dt(Text("Via")),
				Dd(Text(versionMeta.Via)),
			),
		),

		Div(
			Class("control-group"),
			NavButton("Go back", "../.."),
		),
	)
}
