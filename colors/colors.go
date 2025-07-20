package colors

type Colors struct {
	text       string
	background string
}

type ColorIntensity string
type BlackIntensity int

const (
	CI_0 ColorIntensity = "00"
	CI_1 ColorIntensity = "5f"
	CI_2 ColorIntensity = "87"
	CI_3 ColorIntensity = "af"
	CI_4 ColorIntensity = "d7"
	CI_5 ColorIntensity = "ff"
)

const (
	BI_0  BlackIntensity = 0
	BI_1  BlackIntensity = 1
	BI_2  BlackIntensity = 2
	BI_3  BlackIntensity = 3
	BI_4  BlackIntensity = 4
	BI_5  BlackIntensity = 5
	BI_6  BlackIntensity = 6
	BI_7  BlackIntensity = 7
	BI_8  BlackIntensity = 8
	BI_9  BlackIntensity = 9
	BI_10 BlackIntensity = 10
	BI_11 BlackIntensity = 11
	BI_12 BlackIntensity = 12
	BI_13 BlackIntensity = 13
	BI_14 BlackIntensity = 14
	BI_15 BlackIntensity = 15
	BI_16 BlackIntensity = 16
	BI_17 BlackIntensity = 17
	BI_18 BlackIntensity = 18
	BI_19 BlackIntensity = 19
	BI_20 BlackIntensity = 20
	BI_21 BlackIntensity = 21
	BI_22 BlackIntensity = 22
	BI_23 BlackIntensity = 23
)

// NewColors creates a new Colors instance with no values.
func NewColors() *Colors {
	return &Colors{}
}

// SetTextColor sets the text color based on RGB intensity values.
//
// Use CI_0 to CI_5 for each color channel (Red, Green, Blue).
func (c *Colors) SetTextColor(RedIntensity, GreenIntensity, BlueIntensity ColorIntensity) *Colors {
	c.text = "38;5;" + ColorMap[string(RedIntensity)+string(GreenIntensity)+string(BlueIntensity)]
	return c
}

// SetBackgroundColor sets the background color based on RGB intensity values.
//
// Use CI_0 to CI_5 for each color channel (Red, Green, Blue).
func (c *Colors) SetBackgroundColor(RedIntensity, GreenIntensity, BlueIntensity ColorIntensity) *Colors {
	c.background = ";48;5;" + ColorMap[string(RedIntensity)+string(GreenIntensity)+string(BlueIntensity)]
	return c
}

// SetTextGrayScale sets the text color based on grayscale intensity.
//
// Use BI_0 to BI_23 for different grayscale intensities.
func (c *Colors) SetTextGrayScale(BlackIntensity BlackIntensity) *Colors {
	c.text = "38;5;" + greyScale[BlackIntensity]
	return c
}

// SetBackgroundGrayScale sets the background color based on grayscale intensity.
//
// Use BI_0 to BI_10 for different grayscale intensities.
func (c *Colors) SetBackgroundGrayScale(BlackIntensity BlackIntensity) *Colors {
	c.background = ";48;5;" + greyScale[BlackIntensity]
	return c
}

// FormatString formats the string with the set text and background colors.
func (c *Colors) FormatString(s string) string {
	//return `\033[` + c.text + `m` + s + `\033[0m`
	return "\033[" + c.text + c.background + "m" + s + "\033[0m"
}

var ColorMap = map[string]string{
	"000000": "16",
	"00005f": "17",
	"000087": "18",
	"0000af": "19",
	"0000d7": "20",
	"0000ff": "21",
	"005f00": "22",
	"005f5f": "23",
	"005f87": "24",
	"005faf": "25",
	"005fd7": "26",
	"005fff": "27",
	"008700": "28",
	"00875f": "29",
	"008787": "30",
	"0087af": "31",
	"0087d7": "32",
	"0087ff": "33",
	"00af00": "34",
	"00af5f": "35",
	"00af87": "36",
	"00afaf": "37",
	"00afd7": "38",
	"00afff": "39",
	"00d700": "40",
	"00d75f": "41",
	"00d787": "42",
	"00d7af": "43",
	"00d7d7": "44",
	"00d7ff": "45",
	"00ff00": "46",
	"00ff5f": "47",
	"00ff87": "48",
	"00ffaf": "49",
	"00ffd7": "50",
	"00ffff": "51",
	"5f0000": "52",
	"5f005f": "53",
	"5f0087": "54",
	"5f00af": "55",
	"5f00d7": "56",
	"5f00ff": "57",
	"5f5f00": "58",
	"5f5f5f": "59",
	"5f5f87": "60",
	"5f5faf": "61",
	"5f5fd7": "62",
	"5f5fff": "63",
	"5f8700": "64",
	"5f875f": "65",
	"5f8787": "66",
	"5f87af": "67",
	"5f87d7": "68",
	"5f87ff": "69",
	"5faf00": "70",
	"5faf5f": "71",
	"5faf87": "72",
	"5fafaf": "73",
	"5fafd7": "74",
	"5fafff": "75",
	"5fd700": "76",
	"5fd75f": "77",
	"5fd787": "78",
	"5fd7af": "79",
	"5fd7d7": "80",
	"5fd7ff": "81",
	"5fff00": "82",
	"5fff5f": "83",
	"5fff87": "84",
	"5fffaf": "85",
	"5fffd7": "86",
	"5fffff": "87",
	"870000": "88",
	"87005f": "89",
	"870087": "90",
	"8700af": "91",
	"8700d7": "92",
	"8700ff": "93",
	"875f00": "94",
	"875f5f": "95",
	"875f87": "96",
	"875faf": "97",
	"875fd7": "98",
	"875fff": "99",
	"878700": "100",
	"87875f": "101",
	"878787": "102",
	"8787af": "103",
	"8787d7": "104",
	"8787ff": "105",
	"87af00": "106",
	"87af5f": "107",
	"87af87": "108",
	"87afaf": "109",
	"87afd7": "110",
	"87afff": "111",
	"87d700": "112",
	"87d75f": "113",
	"87d787": "114",
	"87d7af": "115",
	"87d7d7": "116",
	"87d7ff": "117",
	"87ff00": "118",
	"87ff5f": "119",
	"87ff87": "120",
	"87ffaf": "121",
	"87fffd": "122",
	"87ffff": "123",
	"af0000": "124",
	"af005f": "125",
	"af0087": "126",
	"af00af": "127",
	"af00d7": "128",
	"af00ff": "129",
	"af5f00": "130",
	"af5f5f": "131",
	"af5f87": "132",
	"af5faf": "133",
	"af5fd7": "134",
	"af5fff": "135",
	"af8700": "136",
	"af875f": "137",
	"af8787": "138",
	"af87af": "139",
	"af87d7": "140",
	"af87ff": "141",
	"afaf00": "142",
	"afaf5f": "143",
	"afaf87": "144",
	"afafaf": "145",
	"afafd7": "146",
	"afafff": "147",
	"afd700": "148",
	"afd75f": "149",
	"afd787": "150",
	"afd7af": "151",
	"afd7d7": "152",
	"afd7ff": "153",
	"afff00": "154",
	"afff5f": "155",
	"afff87": "156",
	"afffaf": "157",
	"afffd7": "158",
	"afffff": "159",
	"d70000": "160",
	"d7005f": "161",
	"d70087": "162",
	"d700af": "163",
	"d700d7": "164",
	"d700ff": "165",
	"d75f00": "166",
	"d75f5f": "167",
	"d75f87": "168",
	"d75faf": "169",
	"d75fd7": "170",
	"d75fff": "171",
	"d78700": "172",
	"d7875f": "173",
	"d78787": "174",
	"d787af": "175",
	"d787d7": "176",
	"d787ff": "177",
	"d7af00": "178",
	"d7af5f": "179",
	"d7af87": "180",
	"d7afaf": "181",
	"d7afd7": "182",
	"d7afff": "183",
	"d7d700": "184",
	"d7d75f": "185",
	"d7d787": "186",
	"d7d7af": "187",
	"d7d7d7": "188",
	"d7d7ff": "189",
	"d7ff00": "190",
	"d7ff5f": "191",
	"d7ff87": "192",
	"d7ffaf": "193",
	"d7ffd7": "194",
	"d7ffff": "195",
	"ff0000": "196",
	"ff005f": "197",
	"ff0087": "198",
	"ff00af": "199",
	"ff00d7": "200",
	"ff00ff": "201",
	"ff5f00": "202",
	"ff5f5f": "203",
	"ff5f87": "204",
	"ff5faf": "205",
	"ff5fd7": "206",
	"ff5fff": "207",
	"ff8700": "208",
	"ff875f": "209",
	"ff8787": "210",
	"ff87af": "211",
	"ff87d7": "212",
	"ff87ff": "213",
	"ffaf00": "214",
	"ffaf5f": "215",
	"ffaf87": "216",
	"ffafaf": "217",
	"ffafd7": "218",
	"ffafff": "219",
	"ffd700": "220",
	"ffd75f": "221",
	"ffd787": "222",
	"ffd7af": "223",
	"ffd7d7": "224",
	"ffd7ff": "225",
	"ffff00": "226",
	"ffff5f": "227",
	"ffff87": "228",
	"ffffaf": "229",
	"ffffd7": "230",
	"ffffff": "231",
}

var greyScale = []string{
	"255",
	"254",
	"253",
	"252",
	"251",
	"250",
	"249",
	"248",
	"247",
	"246",
	"245",
	"244",
	"243",
	"242",
	"241",
	"240",
	"239",
	"238",
	"237",
	"236",
	"235",
	"234",
	"233",
	"232",
}
