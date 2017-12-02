package main

import (
	"regexp"

	"github.com/scott-wilson/dosbot"
)

var baselineRecalibrationMessageResponse = map[*regexp.Regexp]string{
	regexp.MustCompile(`(?i)recite your baseline`):                                                    "And blood-black nothingness began to spin... A system of cells interlinked within cells interlinked within cells interlinked within one stem... And dreadfully distinct against the dark, a tall white fountain played.",
	regexp.MustCompile(`(?i)cells`):                                                                   "Cells.",
	regexp.MustCompile(`(?i)have you been in an institution`):                                         "Cells.",
	regexp.MustCompile(`(?i)do they keep you in a cell`):                                              "Cells.",
	regexp.MustCompile(`(?i)when you're not performing your duties do they keep you in a little box`): "Cells.",
	regexp.MustCompile(`(?i)interlinked`):                                                             "Interlinked.",
	regexp.MustCompile(`(?i)what's it like to hold the hand of someone you love`):                     "Interlinked.",
	regexp.MustCompile(`(?i)did they teach you how to feel finger to finger`):                         "Interlinked.",
	regexp.MustCompile(`(?i)do you long for having your heart interlinked`):                           "Interlinked.",
	regexp.MustCompile(`(?i)do you dream about being interlinked`):                                    "Interlinked.",
	regexp.MustCompile(`(?i)what's it like to hold your child in your arms`):                          "Interlinked.",
	regexp.MustCompile(`(?i)do you feel that there's a part of you that's missing`):                   "Interlinked.",
	regexp.MustCompile(`(?i)within cells interlinked`):                                                "Within cells interlinked.",
	regexp.MustCompile(`(?i)why don't you say within cells interlinked three times`):                  "Within cells interlinked. Within cells interlinked. Within cells interlinked.",
}

func baselineRecalibration(message string) (string, error) {
	for key, value := range baselineRecalibrationMessageResponse {
		if key.FindString(message) != "" {
			return value, nil
		}
	}

	return "", dosbot.ErrEventNotSupportedByAction
}
