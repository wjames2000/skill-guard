package cmd

import (
	"fmt"
	"os"
)

func handleCompletion(args []string) {
	if len(args) < 2 || args[0] != "completion" {
		return
	}

	shell := "bash"
	if len(args) > 1 {
		shell = args[1]
	}

	switch shell {
	case "bash":
		printBashCompletion()
	case "zsh":
		printZshCompletion()
	case "fish":
		printFishCompletion()
	default:
		fmt.Fprintf(os.Stderr, "不支持的 shell: %s (支持: bash, zsh, fish)\n", shell)
		os.Exit(2)
	}
	os.Exit(0)
}

func printBashCompletion() {
	fmt.Println(`_skill_guard_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="-j --json -q --quiet -v --verbose -c --config -r --rules -s --severity -i --ignore --ext-include --ext-exclude --max-size --concurrency --disable-builtin --no-color --output --summary --version -h --help completion"

    case "${prev}" in
        -c|--config|-r|--rules|--output)
            COMPREPLY=($(compgen -f -- "${cur}"))
            return 0
            ;;
        -s|--severity)
            COMPREPLY=($(compgen -W "critical high medium low" -- "${cur}"))
            return 0
            ;;
        *)
            ;;
    esac

    if [[ ${cur} == * ]]; then
        COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
    fi
    return 0
}
complete -F _skill_guard_completion skill-guard`)
}

func printZshCompletion() {
	fmt.Println(`#compdef skill-guard
_skill_guard() {
    local -a opts
    opts=(
        '(-j --json)'{-j,--json}'[JSON format output]'
        '(-q --quiet)'{-q,--quiet}'[Quiet mode]'
        '(-v --verbose)'{-v,--verbose}'[Verbose output]'
        '(-c --config)'{-c,--config}'[Config file path]:file:_files'
        '(-r --rules)'{-r,--rules}'[Custom rules file]:file:_files'
        '(-s --severity)'{-s,--severity}'[Minimum severity level]:level:(critical high medium low)'
        '(-i --ignore)'{-i,--ignore}'[Ignore path pattern]'
        '--ext-include[Extensions to include]'
        '--ext-exclude[Extensions to exclude]'
        '--max-size[Max file size]'
        '--concurrency[Scan concurrency]'
        '--disable-builtin[Disable built-in rules]'
        '--no-color[Disable color output]'
        '--output[Write report to file]:file:_files'
        '--summary[Summary mode]'
        '--sarif[SARIF format output]'
        '--version[Show version]'
        '(-h --help)'{-h,--help}'[Show help]'
        'completion[Generate shell completion]'
    )
    _arguments $opts '*:path:_files -/'
}
_skill_guard`)
}

func printFishCompletion() {
	fmt.Println(`complete -c skill-guard -f
complete -c skill-guard -s j -l json -d "JSON format output"
complete -c skill-guard -s q -l quiet -d "Quiet mode"
complete -c skill-guard -s v -l verbose -d "Verbose output"
complete -c skill-guard -s c -l config -d "Config file path" -r -F
complete -c skill-guard -s r -l rules -d "Custom rules file" -r -F
complete -c skill-guard -s s -l severity -d "Minimum severity level" -r -f -a "critical high medium low"
complete -c skill-guard -s i -l ignore -d "Ignore path pattern" -r
complete -c skill-guard -l ext-include -d "Extensions to include" -r
complete -c skill-guard -l ext-exclude -d "Extensions to exclude" -r
complete -c skill-guard -l max-size -d "Max file size" -r
complete -c skill-guard -l concurrency -d "Scan concurrency" -r
complete -c skill-guard -l disable-builtin -d "Disable built-in rules"
complete -c skill-guard -l no-color -d "Disable color output"
complete -c skill-guard -l output -d "Write report to file" -r -F
complete -c skill-guard -l summary -d "Summary mode"
complete -c skill-guard -l sarif -d "SARIF format output"
complete -c skill-guard -l version -d "Show version"
complete -c skill-guard -s h -l help -d "Show help"
complete -c skill-guard -a completion -d "Generate shell completion"`)
}
