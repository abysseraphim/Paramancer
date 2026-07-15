package internal

import "fmt"

func Detect(baseLine Analysis, analyses []Analysis) []Finding {
	var findings []Finding

	for _, analysis := range analyses {
		if analysis.Error != nil {
			continue
		}

		if analysis.StatusCode != baseLine.StatusCode {

			severity := Medium
			if analysis.StatusCode >= 500 {
				severity = High
			}

			findings = append(findings, Finding{
				Request:  analysis.Request,
				Reason:   "status code changed",
				Severity: severity,
				Evidence: fmt.Sprintf("%d -> %d", baseLine.StatusCode, analysis.StatusCode),
				Mutation: analysis.Request.Mutation,
			})
		}
		if analysis.BodyLength != baseLine.BodyLength {
			findings = append(findings, Finding{
				Request:  analysis.Request,
				Reason:   "body length changed",
				Severity: Low,
				Evidence: fmt.Sprintf("%d bytes -> %d bytes", baseLine.BodyLength, analysis.BodyLength),
				Mutation: analysis.Request.Mutation,
			})
		}
	}

	return findings
}
