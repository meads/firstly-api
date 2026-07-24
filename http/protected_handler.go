package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type dataObject struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Location  string `json:"location"`
	VisitedOn string `json:"visitedOn"`
}

var data = []dataObject{
	{ID: "1", Username: "krolls0", Location: "38.9475624", VisitedOn: "8/5/2025"},
	{ID: "2", Username: "bdagworthy1", Location: "11.5218308", VisitedOn: "1/6/2026"},
	{ID: "3", Username: "emouncher2", Location: "39.81889", VisitedOn: "11/14/2025"},
	{ID: "4", Username: "mivatt3", Location: "10.2593261", VisitedOn: "2/6/2026"},
	{ID: "5", Username: "hdurrett4", Location: "-0.9", VisitedOn: "3/5/2026"},
	{ID: "6", Username: "amerrison5", Location: "-7.1841565", VisitedOn: "12/28/2025"},
	{ID: "7", Username: "gharmeston6", Location: "16.7095", VisitedOn: "7/24/2025"},
	{ID: "8", Username: "hlonglands7", Location: "-7.2457269", VisitedOn: "10/30/2025"},
	{ID: "9", Username: "dbusfield8", Location: "20.775698", VisitedOn: "12/11/2025"},
	{ID: "10", Username: "msinncock9", Location: "20.007968", VisitedOn: "2/16/2026"},
	{ID: "11", Username: "ekinga", Location: "51.8253953", VisitedOn: "9/20/2025"},
	{ID: "12", Username: "hrayburnb", Location: "-6.8332046", VisitedOn: "2/26/2026"},
	{ID: "13", Username: "zlivettc", Location: "31.520703", VisitedOn: "6/25/2026"},
	{ID: "14", Username: "rsexceyd", Location: "48.1933831", VisitedOn: "10/31/2025"},
	{ID: "15", Username: "cagente", Location: "-16.5030766", VisitedOn: "1/8/2026"},
	{ID: "16", Username: "cbewshawf", Location: "15.3098074", VisitedOn: "3/3/2026"},
	{ID: "17", Username: "ahegdeng", Location: "-14.9479605", VisitedOn: "3/13/2026"},
	{ID: "18", Username: "mkingstonh", Location: "38.6037983", VisitedOn: "3/24/2026"},
	{ID: "19", Username: "mgrecei", Location: "56.1597379", VisitedOn: "4/3/2026"},
	{ID: "20", Username: "jdewanej", Location: "48.615982", VisitedOn: "12/6/2025"},
	{ID: "21", Username: "htharmek", Location: "59.4005705", VisitedOn: "1/30/2026"},
	{ID: "22", Username: "delmanl", Location: "14.6593627", VisitedOn: "6/5/2026"},
	{ID: "23", Username: "aflintiffm", Location: "37.5544156", VisitedOn: "9/23/2025"},
	{ID: "24", Username: "kcavelln", Location: "30.081941", VisitedOn: "9/8/2025"},
	{ID: "25", Username: "bpaulo", Location: "-7.4551089", VisitedOn: "11/26/2025"},
	{ID: "26", Username: "dmountlowp", Location: "23.219932", VisitedOn: "11/28/2025"},
	{ID: "27", Username: "shenworthq", Location: "-15.3875259", VisitedOn: "8/31/2025"},
	{ID: "28", Username: "dmeachenr", Location: "-8.6283422", VisitedOn: "9/23/2025"},
	{ID: "29", Username: "ejacmards", Location: "42.891255", VisitedOn: "12/23/2025"},
	{ID: "30", Username: "lclaret", Location: "41.7396715", VisitedOn: "11/21/2025"},
	{ID: "31", Username: "evernu", Location: "43.3487908", VisitedOn: "1/29/2026"},
	{ID: "32", Username: "mabberleyv", Location: "27.9983088", VisitedOn: "3/23/2026"},
	{ID: "33", Username: "hsmalingw", Location: "53.9208518", VisitedOn: "6/9/2026"},
	{ID: "34", Username: "nbeginx", Location: "45.70734", VisitedOn: "1/30/2026"},
	{ID: "35", Username: "kpatryy", Location: "-4.6", VisitedOn: "1/31/2026"},
	{ID: "36", Username: "escanlonz", Location: "59.4084121", VisitedOn: "12/15/2025"},
	{ID: "37", Username: "tvezey10", Location: "40.7465732", VisitedOn: "4/5/2026"},
	{ID: "38", Username: "ddemitris11", Location: "52.8551656", VisitedOn: "3/28/2026"},
	{ID: "39", Username: "kgrebner12", Location: "-23.7047572", VisitedOn: "7/7/2026"},
	{ID: "40", Username: "btideswell13", Location: "47.6451517", VisitedOn: "10/14/2025"},
	{ID: "41", Username: "gfoker14", Location: "-8.6704582", VisitedOn: "6/10/2026"},
	{ID: "42", Username: "lcruickshanks15", Location: "33.7654607", VisitedOn: "5/1/2026"},
	{ID: "43", Username: "ggatus16", Location: "38.7470186", VisitedOn: "2/21/2026"},
	{ID: "44", Username: "grymill17", Location: "24.781681", VisitedOn: "10/26/2025"},
	{ID: "45", Username: "rbunn18", Location: "43.0052216", VisitedOn: "12/1/2025"},
	{ID: "46", Username: "rragbourn19", Location: "-6.831959", VisitedOn: "2/27/2026"},
	{ID: "47", Username: "yjerdon1a", Location: "23.018033", VisitedOn: "9/10/2025"},
	{ID: "48", Username: "jjoly1b", Location: "18.2544504", VisitedOn: "3/4/2026"},
	{ID: "49", Username: "eaccombe1c", Location: "35.2075821", VisitedOn: "3/28/2026"},
	{ID: "50", Username: "ablakely1d", Location: "41.2529921", VisitedOn: "9/22/2025"},
}

func listProtectedHandler(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.JSON(http.StatusOK, data)
}
