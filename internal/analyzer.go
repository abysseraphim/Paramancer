package internal

import (
	"io"
)

func Analyzer(Results []Result) []Analysis {
	var analyses []Analysis

	for _, res := range Results {
		var analysis Analysis

		analysis.Request = res.Request

		if res.Error != nil {
			analysis.Error = res.Error
			analyses = append(analyses, analysis)
			continue
		}
		analysis.StatusCode = res.Response.StatusCode

		body, err := io.ReadAll(res.Response.Body)
		res.Response.Body.Close()
		if err != nil {
			analysis.Error = err
			analyses = append(analyses, analysis)
			continue
		}

		analysis.BodyLength = int64(len(body))
		analysis.ResponseBody = string(body)

		analyses = append(analyses, analysis)
	}

	return analyses
}
