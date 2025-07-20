package themes

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/GenesisBautista/betterminal/colors"
)

type Theme struct {
	Message *colors.Colors `json:"message"`
	Success *colors.Colors `json:"success"`
	Error   *colors.Colors `json:"error"`
	Warning *colors.Colors `json:"warning"`
	Info    *colors.Colors `json:"info"`
	System  *colors.Colors `json:"system"`
	Grid    *colors.Colors `json:"grid"`
}

type themeFromJSON struct {
	Message rgbFromJSON `json:"message"`
	Success rgbFromJSON `json:"success"`
	Error   rgbFromJSON `json:"error"`
	Warning rgbFromJSON `json:"warning"`
	Info    rgbFromJSON `json:"info"`
	System  rgbFromJSON `json:"system"`
	Grid    rgbFromJSON `json:"grid"`
}

type rgbFromJSON struct {
	R *int `json:"R"`
	G *int `json:"G"`
	B *int `json:"B"`
}

// MakeTheme creates a default theme with predefined colors.
func MakeTheme() *Theme {
	return &Theme{
		Message: colors.NewColors().SetTextGrayScale(colors.BI_2),
		Success: colors.NewColors().SetTextColor(colors.CI_2, colors.CI_5, colors.CI_2),
		Error:   colors.NewColors().SetTextColor(colors.CI_5, colors.CI_2, colors.CI_2),
		Warning: colors.NewColors().SetTextColor(colors.CI_5, colors.CI_5, colors.CI_2),
		Info:    colors.NewColors().SetTextColor(colors.CI_2, colors.CI_2, colors.CI_5),
		System:  colors.NewColors().SetTextGrayScale(colors.BI_10),
		Grid:    colors.NewColors().SetTextGrayScale(colors.BI_20),
	}
}

// MakeThemeFromJSON creates a theme from a JSON file.
//
// If the JSON file is not found or cannot be parsed, it returns a default theme.
//
// JSON format should be:
//
//	{
//	  "message": {"R": 0, "G": 5, "B": 0},
//	  "success": {"R": 0, "G": 5, "B": 0},
//	  "error": {"R": 5, "G": 0, "B": 0},
//	  "warning": {"R": 5, "G": 5, "B": 0},
//	  "info": {"R": 0, "G": 0, "B": 5},
//	  "system": {"G": 23},
//	  "grid": {"G": 2}
//	}
//
// RGB values range from 0 to 5 for each channel unless you are using gray scale then the range is 0-23.
func MakeThemeFromJSON(jsonPath string) (*Theme, error) {
	bytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return MakeTheme(), err
	}

	t := &themeFromJSON{}
	if err := json.Unmarshal(bytes, t); err != nil {
		return MakeTheme(), err
	}

	return parseThemeFromJsonToTheme(t)
}

func parseThemeFromJsonToTheme(t *themeFromJSON) (*Theme, error) {
	theme := &Theme{}
	var err error

	if theme.Message, err = parseRgb(t.Message); err != nil {
		return nil, err
	}

	if theme.Success, err = parseRgb(t.Success); err != nil {
		return nil, err
	}

	if theme.Error, err = parseRgb(t.Error); err != nil {
		return nil, err
	}

	if theme.Warning, err = parseRgb(t.Warning); err != nil {
		return nil, err
	}

	if theme.Info, err = parseRgb(t.Info); err != nil {
		return nil, err
	}

	if theme.System, err = parseRgb(t.System); err != nil {
		return nil, err
	}

	if theme.Grid, err = parseRgb(t.Grid); err != nil {
		return nil, err
	}

	return theme, nil
}

func parseRgb(rgb rgbFromJSON) (*colors.Colors, error) {
	c := colors.NewColors()

	if rgb.R == nil && rgb.G == nil {
		return nil, errors.New("could not find colors please specify RGB values or just G for gray scale")
	} else if rgb.R == nil && rgb.G != nil {
		if *rgb.G >= 0 && *rgb.G <= 23 {
			return c.SetTextGrayScale(colors.BlackIntensity(*rgb.G)), nil
		} else {
			return nil, errors.New("gray scale out of bounds(0-23)")
		}
	} else {
		if (*rgb.R >= 0 && *rgb.R <= 5) && (*rgb.G >= 0 && *rgb.G <= 5) && (*rgb.B >= 0 && *rgb.B <= 5) {
			ci := []string{"00", "5f", "87", "af", "d7", "ff"}
			return c.SetTextColor(colors.ColorIntensity(ci[*rgb.R]), colors.ColorIntensity(ci[*rgb.G]), colors.ColorIntensity(ci[*rgb.B])), nil
		} else {
			return nil, errors.New("RGB out of bounds(0-5)")
		}
	}
}

func (t *Theme) PrintlnMessage(msg string) {
	fmt.Println(t.Message.FormatString(msg))
}

func (t *Theme) PrintlnSuccess(msg string) {
	fmt.Println(t.Success.FormatString(msg))
}

func (t *Theme) PrintlnError(msg string) {
	fmt.Println(t.Error.FormatString(msg))
}

func (t *Theme) PrintlnWarning(msg string) {
	fmt.Println(t.Warning.FormatString(msg))
}

func (t *Theme) PrintlnInfo(msg string) {
	fmt.Println(t.Info.FormatString(msg))
}
