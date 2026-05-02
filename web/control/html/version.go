package html

import (
	"strconv"
	"time"

	"github.com/LMBishop/scrapbook/pkg/config"
	. "github.com/LMBishop/scrapbook/web/skeleton"
	"github.com/dustin/go-humanize"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func VersionPage(versionMeta config.VersionMeta) Node {
	createDate := time.Unix(int64(versionMeta.Created), 0)

	return Page("Version "+versionMeta.Hash[0:8],
		H1(Text("Version "+versionMeta.Hash[0:8])),

		FieldSet(
			Legend(Text("Version actions")),

			Div(
				Class("control-group"),

				NavButton("Download version", "#"),
				NavButton("Delete version", "delete"),
			),
		),

		P(
			Dl(
				Dt(Text("Hash")),
				Dd(Text(versionMeta.Hash)),

				Dt(Text("Number of files")),
				Dd(Text(strconv.Itoa(int(versionMeta.Files)))),

				Dt(Text("Size on disk")),
				Dd(Text(humanize.IBytes(uint64(versionMeta.Size)))),

				Dt(Text("Created at")),
				Dd(Text(createDate.Format(time.UnixDate))),

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
