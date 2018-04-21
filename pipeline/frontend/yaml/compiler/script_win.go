package compiler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
)

func generateScriptWindows(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScriptWin,
			escaped,
			command,
		))
	}
	script := fmt.Sprintf(
		setupScriptWin,
		buf.String(),
	)
	return base64.StdEncoding.EncodeToString([]byte(script))
}

// TODO empty CI_NETRC_MACHINE check
const setupScriptWin = `
$ErrorActionPreference = 'Stop';
if ($Env:CI_NETRC_MACHINE) {
$netrc=[string]::Format("{0}\_netrc",$Env:USERPROFILE);
"machine $Env:CI_NETRC_MACHINE" >> $netrc;
"login $Env:CI_NETRC_USERNAME" >> $netrc;
"password $Env:CI_NETRC_PASSWORD" >> $netrc;
};
[Environment]::SetEnvironmentVariable("CI_NETRC_PASSWORD",$null);
[Environment]::SetEnvironmentVariable("CI_SCRIPT",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_USERNAME",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_PASSWORD",$null)
&cmd /c "mkdir c:\root"
%s
`

// traceScript is a helper script that is added to the build script
// to trace a command.
const traceScriptWin = `
Write-Output ('+ %s');
&cmd /c "%s";
`
