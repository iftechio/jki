package completion

const (
	bashCompletionFunc = `# call img $1,
__jki_debug_out()
{
    local cmd="$1"
    __jki_debug "${FUNCNAME[1]}: get completion by ${cmd}"
    eval "${cmd} 2>/dev/null"
}
__jki_override_flag_list=(--kubeconfig --context --namespace -n)
__jki_override_flags()
{
    local ${__jki_override_flag_list[*]##*-} two_word_of of var
    for w in "${words[@]}"; do
        if [ -n "${two_word_of}" ]; then
            eval "${two_word_of##*-}=\"${two_word_of}=\${w}\""
            two_word_of=
            continue
        fi
        for of in "${__jki_override_flag_list[@]}"; do
            case "${w}" in
                ${of}=*)
                    eval "${of##*-}=\"${w}\""
                    ;;
                ${of})
                    two_word_of="${of}"
                    ;;
            esac
        done
    done
    for var in "${__jki_override_flag_list[@]##*-}"; do
        if eval "test -n \"\$${var}\""; then
            eval "echo -n \${${var}}' '"
        fi
    done
}

# $1 has to be "contexts", "servers", "users"
__jki_parse_config()
{
    local template jki_out
    template="{{ range .$1  }}{{ .name }} {{ end }}"
    if jki_out=$(__jki_debug_out "kubectl config $(__jki_override_flags) -o template --template=\"${template}\" view"); then
        COMPREPLY=( $( compgen -W "${jki_out[*]}" -- "$cur" ) )
    fi
}
__jki_config_get_contexts()
{
    __jki_parse_config "contexts"
}

# $1 is the name of resource (required)
# $2 is template string for kubectl get (optional)
__jki_parse_get()
{
    local template
    template="${2:-"{{ range .items  }}{{ .metadata.name }} {{ end }}"}"
    local jki_out
    if jki_out=$(__jki_debug_out "kubectl get $(__jki_override_flags) -o template --template=\"${template}\" \"$1\""); then
        COMPREPLY+=( $( compgen -W "${jki_out[*]}" -- "$cur" ) )
    fi
}
__jki_get_resource()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
      local jki_out
      if jki_out=$(__jki_debug_out "kubectl api-resources $(__kt_override_flags) -o name --cached --request-timeout=5s --verbs=get"); then
          COMPREPLY=( $( compgen -W "${jki_out[*]}" -- "$cur" ) )
          return 0
      fi
      return 1
    fi
    __jki_parse_get "${nouns[${#nouns[@]} -1]}"
}
__jki_get_resource_namespace()
{
    __jki_parse_get "namespace"
}

__jki_config_get_registries()
{
    local template jki_out
    if jki_out=$(__jki_debug_out "img config get-registries"); then
        COMPREPLY=( $( compgen -W "${jki_out[*]}" -- "$cur" ) )
    fi
}

__jki_abort() {
    return 1
}

`
)

var (
	bashCompletionFlags = map[string]string{
		"namespace": "__jki_get_resource_namespace",
		"context":   "__jki_config_get_contexts",
		"registry":  "__jki_config_get_registries",

		"kubeconfig": "__jki_abort",
		"jkiconfig":  "__jki_abort",
	}
)
