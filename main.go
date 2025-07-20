package main

import (
	"fmt"

	"github.com/GenesisBautista/betterminal/colors"
	"github.com/GenesisBautista/betterminal/themes"
)

func main() {
	colors_test()
	greyScale_test()
	themes_test()
	themesJson_test()
}

func themes_test() {
	t := themes.MakeTheme()

	t.PrintlnMessage("This is a message.")
	t.PrintlnSuccess("This is a success message.")
	t.PrintlnError("This is an error message.")
	t.PrintlnWarning("This is a warning message.")
	t.PrintlnInfo("This is an info message.")
}

func themesJson_test() {
	t, err := themes.MakeThemeFromJSON("theme.json")
	if err != nil {
		fmt.Println("Error loading theme:", err)
		return
	}

	t.PrintlnMessage("This is a message from JSON theme.")
	t.PrintlnSuccess("This is a success message from JSON theme.")
	t.PrintlnError("This is an error message from JSON theme.")
	t.PrintlnWarning("This is a warning message from JSON theme.")
	t.PrintlnInfo("This is an info message from JSON theme.")
}

func colors_test() {
	c := colors.Colors{}

	c.SetTextColor(colors.CI_5, colors.CI_0, colors.CI_0)
	fmt.Println(c.FormatString("Hello, World!"))
	c.SetBackgroundColor(colors.CI_0, colors.CI_5, colors.CI_5)
	fmt.Println(c.FormatString("Hello, World! with background"))
}

func greyScale_test() {
	c := colors.Colors{}

	c.SetTextGrayScale(colors.BI_5)
	fmt.Println(c.FormatString("Hello, World! in grayscale"))

	c.SetBackgroundGrayScale(colors.BI_0)
	fmt.Println(c.FormatString("Intenity 0 background"))
	c.SetBackgroundGrayScale(colors.BI_1)
	fmt.Println(c.FormatString("Intenity 1 background"))
	c.SetBackgroundGrayScale(colors.BI_2)
	fmt.Println(c.FormatString("Intenity 2 background"))
	c.SetBackgroundGrayScale(colors.BI_3)
	fmt.Println(c.FormatString("Intenity 3 background"))
	c.SetBackgroundGrayScale(colors.BI_4)
	fmt.Println(c.FormatString("Intenity 4 background"))
	c.SetBackgroundGrayScale(colors.BI_5)
	fmt.Println(c.FormatString("Intenity 5 background"))
	c.SetBackgroundGrayScale(colors.BI_6)
	fmt.Println(c.FormatString("Intenity 6 background"))
	c.SetBackgroundGrayScale(colors.BI_7)
	fmt.Println(c.FormatString("Intenity 7 background"))
	c.SetBackgroundGrayScale(colors.BI_8)
	fmt.Println(c.FormatString("Intenity 8 background"))
	c.SetBackgroundGrayScale(colors.BI_9)
	fmt.Println(c.FormatString("Intenity 9 background"))
	c.SetBackgroundGrayScale(colors.BI_10)
	fmt.Println(c.FormatString("Intenity 10 background"))
	c.SetBackgroundGrayScale(colors.BI_11)
	fmt.Println(c.FormatString("Intenity 11 background"))
	c.SetBackgroundGrayScale(colors.BI_12)
	fmt.Println(c.FormatString("Intenity 12 background"))
	c.SetBackgroundGrayScale(colors.BI_13)
	fmt.Println(c.FormatString("Intenity 13 background"))
	c.SetBackgroundGrayScale(colors.BI_14)
	fmt.Println(c.FormatString("Intenity 14 background"))
	c.SetBackgroundGrayScale(colors.BI_15)
	fmt.Println(c.FormatString("Intenity 15 background"))
	c.SetBackgroundGrayScale(colors.BI_16)
	fmt.Println(c.FormatString("Intenity 16 background"))
	c.SetBackgroundGrayScale(colors.BI_17)
	fmt.Println(c.FormatString("Intenity 17 background"))
	c.SetBackgroundGrayScale(colors.BI_18)
	fmt.Println(c.FormatString("Intenity 18 background"))
	c.SetBackgroundGrayScale(colors.BI_19)
	fmt.Println(c.FormatString("Intenity 19 background"))
	c.SetBackgroundGrayScale(colors.BI_20)
	fmt.Println(c.FormatString("Intenity 20 background"))
	c.SetBackgroundGrayScale(colors.BI_21)
	fmt.Println(c.FormatString("Intenity 21 background"))
	c.SetBackgroundGrayScale(colors.BI_22)
	fmt.Println(c.FormatString("Intenity 22 background"))
	c.SetBackgroundGrayScale(colors.BI_23)
	fmt.Println(c.FormatString("Intenity 23 background"))
}
