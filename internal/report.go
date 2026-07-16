package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Reporter(request Request, baseline Analysis, findings []Finding) error {
	var report Report

	report.Target = request.URL
	report.Baseline = baseline

	for _, finding := range findings {
		var f FinalResult

		f.Severity = finding.Severity
		f.Reason = finding.Reason
		f.Evidence = finding.Evidence

		f.Parameter = finding.Mutation.Parameter
		f.ParamLocation = finding.Mutation.Location
		f.ParamTargetedPart = finding.Mutation.Target
		f.MutationStrategy = finding.Mutation.Strategy
		orig := ""
		for _, p := range request.Parameters {
			if p.Name == finding.Mutation.Parameter && p.Location == finding.Mutation.Location {
				if finding.Mutation.Target == NameTarget {
					orig = p.Name
				} else {
					orig = p.Value
				}
				break
			}
		}
		switch finding.Mutation.Strategy {
		case Append:
			f.Payload = orig + finding.Mutation.Payload
		case Prepend:
			f.Payload = finding.Mutation.Payload + orig
		case Replace:
			f.Payload = finding.Mutation.Payload
		}

		report.Findings = append(report.Findings, f)
	}

	PrintCli(report)
	err := SaveJSON(report)
	if err != nil {
		return err
	}

	return nil
}

func PrintCli(rep Report) {
	var buildCli strings.Builder
	buildCli.WriteString(" ====================\n")
	buildCli.WriteString("      Paramancer     \n")
	buildCli.WriteString(" ====================\n\n")

	buildCli.WriteString(" Target\n")
	buildCli.WriteString(fmt.Sprintf(" %v", rep.Target))
	buildCli.WriteString("\n\n")

	buildCli.WriteString(" Baseline info:\n")
	buildCli.WriteString(" Status Code:\n")
	buildCli.WriteString(fmt.Sprintf(" %v", strconv.Itoa(rep.Baseline.StatusCode)))
	buildCli.WriteString("\n\n")

	buildCli.WriteString(" Findings:\n")
	buildCli.WriteString("\n")

	for _, finding := range rep.Findings {
		buildCli.WriteString(" [+] severity: ")
		buildCli.WriteString(string(finding.Severity))
		buildCli.WriteString("\n")

		buildCli.WriteString(" [+] reason: ")
		buildCli.WriteString(finding.Reason)
		buildCli.WriteString("\n")

		buildCli.WriteString(" [+] evidence: ")
		buildCli.WriteString(finding.Evidence)
		buildCli.WriteString("\n")

		buildCli.WriteString(" [+] mutation:\n")
		buildCli.WriteString("    [++] parameter: ")
		buildCli.WriteString(finding.Parameter)
		buildCli.WriteString("\n")

		buildCli.WriteString("    [++] location: ")
		buildCli.WriteString(string(finding.ParamLocation))
		buildCli.WriteString("\n")

		buildCli.WriteString("    [++] target: ")
		buildCli.WriteString(string(finding.ParamTargetedPart))
		buildCli.WriteString("\n")

		buildCli.WriteString("    [++] strategy: ")
		buildCli.WriteString(string(finding.MutationStrategy))
		buildCli.WriteString("\n")

		buildCli.WriteString("    [++] payload: ")
		buildCli.WriteString(finding.Payload)

		buildCli.WriteString("\n\n-------------------\n\n")
	}

	fmt.Print(buildCli.String())
}

func SaveJSON(rep Report) error {
	rep.Baseline.ResponseBody = "empty by choice"
	rep.Baseline.Request.Mutation = Mutation{}

	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		fmt.Printf("\033[31m  [!] failed to create JSON structure\033[0m")
		return err
	}
	data = bytes.ReplaceAll(data, []byte(`\u003e`), []byte(`>`))

	var filename string

	targetURL, err := url.Parse(rep.Target)
	if err != nil {
		filename = fmt.Sprintf("%v.output.json", rep.Target)
	} else {
		filename = targetURL.Host + ".output.json"
	}

	filename = strings.ReplaceAll(filename, ":", "-")
	filename = strings.ReplaceAll(filename, "/", "_")

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("\033[31m  [!] failed to create JSON file\033[0m")
		return err
	}
	fmt.Println()
	fmt.Println("\033[32m  [*] results also saved to: \033[0m", filename)
	fmt.Println()

	return nil
}
