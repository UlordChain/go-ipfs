_udfs_comp()
{
    COMPREPLY=( $(compgen -W "$1" -- ${word}) )
    if [[ ${#COMPREPLY[@]} == 1 && ${COMPREPLY[0]} == "--"*"=" ]] ; then
        # If there's only one option, with =, then discard space
        compopt -o nospace
    fi
}

_udfs_help_only()
{
    _udfs_comp "--help"
}

_udfs_add()
{
    if [[ "${prev}" == "--chunker" ]] ; then
        _udfs_comp "placeholder1 placeholder2 placeholder3" # TODO: a) Give real options, b) Solve autocomplete bug for "="
    elif [ "${prev}" == "--pin" ] ; then
        _udfs_comp "true false"
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive --quiet --silent --progress --trickle --only-hash --wrap-with-directory --hidden --chunker= --pin= --raw-leaves --help "
    else
        _udfs_filesystem_complete
    fi
}

_udfs_bitswap()
{
    udfs_comp "ledger stat unwant wantlist --help"
}

_udfs_bitswap_ledger()
{
    _udfs_help_only
}

_udfs_bitswap_stat()
{
    _udfs_help_only
}

_udfs_bitswap_unwant()
{
    _udfs_help_only
}

_udfs_bitswap_wantlist()
{
    udfs_comp "--peer= --help"
}

_udfs_bitswap_unwant()
{
    _udfs_help_only
}

_udfs_block()
{
    _udfs_comp "get put rm stat --help"
}

_udfs_block_get()
{
    _udfs_hash_complete
}

_udfs_block_put()
{
    if [ "${prev}" == "--format" ] ; then
        _udfs_comp "v0 placeholder2 placeholder3" # TODO: a) Give real options, b) Solve autocomplete bug for "="
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--format= --help"
    else
        _udfs_filesystem_complete
    fi
}

_udfs_block_rm()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--force --quiet --help"
    else
        _udfs_hash_complete
    fi
}

_udfs_block_stat()
{
    _udfs_hash_complete
}

_udfs_bootstrap()
{
    _udfs_comp "add list rm --help"
}

_udfs_bootstrap_add()
{
    _udfs_comp "default --help"
}

_udfs_bootstrap_list()
{
    _udfs_help_only
}

_udfs_bootstrap_rm()
{
    _udfs_comp "all --help"
}

_udfs_cat()
{
    if [[ ${prev} == */* ]] ; then
        COMPREPLY=() # Only one argument allowed
    elif [[ ${word} == */* ]] ; then
        _udfs_hash_complete
    else
        _udfs_pinned_complete
    fi
}

_udfs_commands()
{
    _udfs_comp "--flags --help"
}

_udfs_config()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--bool --json"
    elif [[ ${prev} == *.* ]] ; then
        COMPREPLY=() # Only one subheader of the config can be shown or edited.
    else
        _udfs_comp "show edit replace"
    fi
}

_udfs_config_edit()
{
    _udfs_help_only
}

_udfs_config_replace()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--help"
    else
        _udfs_filesystem_complete
    fi
}

_udfs_config_show()
{
    _udfs_help_only
}

_udfs_daemon()
{
    if [[ ${prev} == "--routing" ]] ; then
        _udfs_comp "dht dhtclient none" # TODO: Solve autocomplete bug for "="
    elif [[ ${prev} == "--mount-udfs" ]] || [[ ${prev} == "--mount-ipns" ]] || [[ ${prev} == "=" ]]; then
        _udfs_filesystem_complete
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--init --routing= --mount --writable --mount-udfs= \
            --mount-ipns= --unrestricted-api --disable-transport-encryption \
            -- enable-gc --manage-fdlimit --offline --migrate --help"
    fi
}

_udfs_dag()
{
    _udfs_comp "get put --help"
}

_udfs_dag_get()
{
    _udfs_help_only
}

_udfs_dag_put()
{
    if [[ ${prev} == "--format" ]] ; then
        _udfs_comp "cbor placeholder1" # TODO: a) Which format more then cbor is valid? b) Solve autocomplete bug for "="
    elif [[ ${prev} == "--input-enc" ]] ; then
        _udfs_comp "json placeholder1" # TODO: a) Which format more then json is valid? b) Solve autocomplete bug for "="
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--format= --input-enc= --help"
    else
        _udfs_filesystem_complete
    fi
}

_udfs_dht()
{
    _udfs_comp "findpeer findprovs get provide put query --help"
}

_udfs_dht_findpeer()
{
    _udfs_comp "--verbose --help"
}

_udfs_dht_findprovs()
{
    _udfs_comp "--verbose --help"
}

_udfs_dht_get()
{
    _udfs_comp "--verbose --help"
}

_udfs_dht_provide()
{
    _udfs_comp "--recursive --verbose --help"
}

_udfs_dht_put()
{
    _udfs_comp "--verbose --help"
}

_udfs_dht_query()
{
    _udfs_comp "--verbose --help"
}

_udfs_diag()
{
    _udfs_comp "sys cmds net --help"
}

_udfs_diag_cmds()
{
    if [[ ${prev} == "clear" ]] ; then
        return 0
    elif [[ ${prev} =~ ^-?[0-9]+$ ]] ; then
        _udfs_comp "ns us µs ms s m h" # TODO: Trigger with out space, eg. "udfs diag set-time 10ns" not "... set-time 10 ns"
    elif [[ ${prev} == "set-time" ]] ; then
        _udfs_help_only
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--verbose --help"
    else
        _udfs_comp "clear set-time"
    fi
}

_udfs_diag_sys()
{
    _udfs_help_only
}

_udfs_diag_net()
{
    if [[ ${prev} == "--vis" ]] ; then
        _udfs_comp "d3 dot text" # TODO: Solve autocomplete bug for "="
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--timeout= --vis= --help"
    fi
}

_udfs_dns()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive --help"
    fi
}

_udfs_files()
{
    _udfs_comp "mv rm flush read write cp ls mkdir stat"
}

_udfs_files_mv()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive --flush"
    elif [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_files_rm()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive --flush"
    elif [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}
_udfs_files_flush()
{
    if [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_files_read()
{
    if [[ ${prev} == "--count" ]] || [[ ${prev} == "--offset" ]] ; then
        COMPREPLY=() # Numbers, just keep it empty
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--offset --count --help"
    elif [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_files_write()
{
    if [[ ${prev} == "--count" ]] || [[ ${prev} == "--offset" ]] ; then # Dirty check
        COMPREPLY=() # Numbers, just keep it empty
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--offset --count --create --truncate --help"
    elif [[ ${prev} == /* ]] ; then
        _udfs_filesystem_complete
    elif [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_files_cp()
{
    if [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_files_ls()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "-l --help"
    elif [[ ${prev} == /* ]] ; then
        COMPREPLY=() # Path exist
    elif [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_files_mkdir()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--parents --help"

    elif [[ ${prev} == /* ]] ; then
        COMPREPLY=() # Path exist
    elif [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_files_stat()
{
    if [[ ${prev} == /* ]] ; then
        COMPREPLY=() # Path exist
    elif [[ ${word} == /* ]] ; then
        _udfs_files_complete
    else
        COMPREPLY=( / )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_file()
{
    if [[ ${prev} == "ls" ]] ; then
        _udfs_hash_complete
    else
        _udfs_comp "ls --help"
    fi
}

_udfs_file_ls()
{
    _udfs_help_only
}

_udfs_get()
{
    if [ "${prev}" == "--output" ] ; then
        compopt -o default # Re-enable default file read
        COMPREPLY=()
    elif [ "${prev}" == "--compression-level" ] ; then
        _udfs_comp "-1 1 2 3 4 5 6 7 8 9" # TODO: Solve autocomplete bug for "="
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--output= --archive --compress --compression-level= --help"
    else
        _udfs_hash_complete
    fi
}

_udfs_id()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--format= --help"
    fi
}

_udfs_init()
{
    _udfs_comp "--bits --force --empty-repo --help"
}

_udfs_log()
{
    _udfs_comp "level ls tail --help"
}

_udfs_log_level()
{
    # TODO: auto-complete subsystem and level
    _udfs_help_only
}

_udfs_log_ls()
{
    _udfs_help_only
}

_udfs_log_tail()
{
    _udfs_help_only
}

_udfs_ls()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--headers --resolve-type=false --help"
    else
        _udfs_hash_complete
    fi
}

_udfs_mount()
{
    if [[ ${prev} == "--udfs-path" ]] || [[ ${prev} == "--ipns-path" ]] || [[ ${prev} == "=" ]] ; then
        _udfs_filesystem_complete
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--udfs-path= --ipns-path= --help"
    fi
}

_udfs_name()
{
    _udfs_comp "publish resolve --help"
}

_udfs_name_publish()
{
    if [[ ${prev} == "--lifetime" ]] || [[ ${prev} == "--ttl" ]] ; then
        COMPREPLY=() # Accept only numbers
    elif [[ ${prev} =~ ^-?[0-9]+$ ]] ; then
        _udfs_comp "ns us µs ms s m h" # TODO: Trigger without space, eg. "udfs diag set-time 10ns" not "... set-time 10 ns"
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--resolve --lifetime --ttl --help"
    elif [[ ${word} == */ ]]; then
        _udfs_hash_complete
    else
        _udfs_pinned_complete
    fi
}

_udfs_name_resolve()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive --nocache --help"
    fi
}

_udfs_object()
{
    _udfs_comp "data diff get links new patch put stat --help"
}

_udfs_object_data()
{
    _udfs_hash_complete
}

_udfs_object_diff()
{
  if [[ ${word} == -* ]] ; then
      _udfs_comp "--verbose --help"
  else
      _udfs_hash_complete
  fi
}


_udfs_object_get()
{
    if [ "${prev}" == "--encoding" ] ; then
        _udfs_comp "protobuf json xml"
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--encoding --help"
    else
        _udfs_hash_complete
    fi
}

_udfs_object_links()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--headers --help"
    else
        _udfs_hash_complete
    fi
}

_udfs_object_new()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--help"
    else
        _udfs_comp "unixfs-dir"
    fi
}

_udfs_object_patch()
{
    if [[ -n "${COMP_WORDS[3]}" ]] ; then # Root merkledag object exist
        case "${COMP_WORDS[4]}" in
        append-data)
            _udfs_help_only
            ;;
        add-link)
            if [[ ${word} == -* ]] && [[ ${prev} == "add-link" ]] ; then # Dirty check
                _udfs_comp "--create"
            #else
                # TODO: Hash path autocomplete. This is tricky, can be hash or a name.
            fi
            ;;
        rm-link)
            _udfs_hash_complete
            ;;
        set-data)
            _udfs_filesystem_complete
            ;;
        *)
            _udfs_comp "append-data add-link rm-link set-data"
            ;;
        esac
    else
        _udfs_hash_complete
    fi
}

_udfs_object_put()
{
    if [ "${prev}" == "--inputenc" ] ; then
        _udfs_comp "protobuf json"
    elif [ "${prev}" == "--datafieldenc" ] ; then
        _udfs_comp "text base64"
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--inputenc --datafieldenc --help"
    else
        _udfs_hash_complete
    fi
}

_udfs_object_stat()
{
    _udfs_hash_complete
}

_udfs_pin()
{
    _udfs_comp "rm ls add --help"
}

_udfs_pin_add()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive=  --help"
    elif [[ ${word} == */ ]] && [[ ${word} != "/udfs/" ]] ; then
        _udfs_hash_complete
    fi
}

_udfs_pin_ls()
{
    if [[ ${prev} == "--type" ]] || [[ ${prev} == "-t" ]] ; then
        _udfs_comp "direct indirect recursive all" # TODO: Solve autocomplete bug for
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--count --quiet --type= --help"
    elif [[ ${word} == */ ]] && [[ ${word} != "/udfs/" ]] ; then
        _udfs_hash_complete
    fi
}

_udfs_pin_rm()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive  --help"
    elif [[ ${word} == */ ]] && [[ ${word} != "/udfs/" ]] ; then
        COMPREPLY=() # TODO: _udfs_hash_complete() + List local pinned hashes as default?
    fi
}

_udfs_ping()
{
    _udfs_comp "--count=  --help"
}

_udfs_pubsub()
{
    _udfs_comp "ls peers pub sub --help"
}

_udfs_pubsub_ls()
{
    _udfs_help_only
}

_udfs_pubsub_peers()
{
    _udfs_help_only
}

_udfs_pubsub_pub()
{
    _udfs_help_only
}

_udfs_pubsub_sub()
{
    _udfs_comp "--discover --help"
}

_udfs_refs()
{
    if [ "${prev}" == "--format" ] ; then
        _udfs_comp "src dst linkname"
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "local --format= --edges --unique --recursive --help"
    #else
        # TODO: Use "udfs ref" and combine it with autocomplete, see _udfs_hash_complete
    fi
}

_udfs_refs_local()
{
    _udfs_help_only
}

_udfs_repo()
{
    _udfs_comp "fsck gc stat verify version --help"
}

_udfs_repo_version()
{
    _udfs_comp "--quiet --help"
}

_udfs_repo_verify()
{
    _udfs_help_only
}

_udfs_repo_gc()
{
    _udfs_comp "--quiet --help"
}

_udfs_repo_stat()
{
    _udfs_comp "--human --help"
}

_udfs_repo_fsck()
{
    _udfs_help_only
}

_udfs_resolve()
{
    if [[ ${word} == /udfs/* ]] ; then
        _udfs_hash_complete
    elif [[ ${word} == /ipns/* ]] ; then
        COMPREPLY=() # Can't autocomplete ipns
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--recursive --help"
    else
        opts="/ipns/ /udfs/"
        COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
        [[ $COMPREPLY = */ ]] && compopt -o nospace
    fi
}

_udfs_stats()
{
    _udfs_comp "bitswap bw repo --help"
}

_udfs_stats_bitswap()
{
    _udfs_help_only
}

_udfs_stats_bw()
{
    # TODO: Which protocol is valid?
    _udfs_comp "--peer= --proto= --poll --interval= --help"
}

_udfs_stats_repo()
{
    _udfs_comp "--human= --help"
}

_udfs_swarm()
{
    _udfs_comp "addrs connect disconnect filters peers --help"
}

_udfs_swarm_addrs()
{
    _udfs_comp "local --help"
}

_udfs_swarm_addrs_local()
{
    _udfs_comp "--id --help"
}

_udfs_swarm_connect()
{
    _udfs_multiaddr_complete
}

_udfs_swarm_disconnect()
{
    local OLDIFS="$IFS" ; local IFS=$'\n' # Change divider for iterator one line below
    opts=$(for x in `udfs swarm peers`; do echo ${x} ; done)
    IFS="$OLDIFS" # Reset divider to space, ' '
    COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
    [[ $COMPREPLY = */ ]] && compopt -o nospace -o filenames
}

_udfs_swarm_filters()
{
    if [[ ${prev} == "add" ]] || [[ ${prev} == "rm" ]]; then
        _udfs_multiaddr_complete
    else
        _udfs_comp "add rm --help"
    fi
}

_udfs_swarm_filters_add()
{
    _udfs_help_only
}

_udfs_swarm_filters_rm()
{
    _udfs_help_only
}

_udfs_swarm_peers()
{
    _udfs_help_only
}

_udfs_tar()
{
    _udfs_comp "add cat --help"
}

_udfs_tar_add()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--help"
    else
        _udfs_filesystem_complete
    fi
}

_udfs_tar_cat()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--help"
    else
        _udfs_filesystem_complete
    fi
}

_udfs_update()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--version" # TODO: How does "--verbose" option work?
    else
        _udfs_comp "versions version install stash revert fetch"
    fi
}

_udfs_update_install()
{
    if   [[ ${prev} == v*.*.* ]] ; then
        COMPREPLY=()
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--version"
    else
        local OLDIFS="$IFS" ; local IFS=$'\n' # Change divider for iterator one line below
        opts=$(for x in `udfs update versions`; do echo ${x} ; done)
        IFS="$OLDIFS" # Reset divider to space, ' '
        COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
    fi
}

_udfs_update_stash()
{
    if [[ ${word} == -* ]] ; then
        _udfs_comp "--tag --help"
    fi
}
_udfs_update_fetch()
{
    if [[ ${prev} == "--output" ]] ; then
        _udfs_filesystem_complete
    elif [[ ${word} == -* ]] ; then
        _udfs_comp "--output --help"
    fi
}

_udfs_version()
{
    _udfs_comp "--number --commit --repo"
}

_udfs_hash_complete()
{
    local lastDir=${word%/*}/
    echo "LastDir: ${lastDir}" >> ~/Downloads/debug-udfs.txt
    local OLDIFS="$IFS" ; local IFS=$'\n' # Change divider for iterator one line below
    opts=$(for x in `udfs file ls ${lastDir}`; do echo ${lastDir}${x}/ ; done) # TODO: Implement "udfs file ls -F" to get rid of frontslash after files. This take long time to run first time on a new shell.
    echo "Options: ${opts}" >> ~/Downloads/debug-udfs.txt
    IFS="$OLDIFS" # Reset divider to space, ' '
    echo "Current: ${word}" >> ~/Downloads/debug-udfs.txt
    COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
    echo "Suggestion: ${COMPREPLY}" >> ~/Downloads/debug-udfs.txt
    [[ $COMPREPLY = */ ]] && compopt -o nospace -o filenames # Removing whitespace after output & handle output as filenames. (Only printing the latest folder of files.)
    return 0
}

_udfs_files_complete()
{
    local lastDir=${word%/*}/
    local OLDIFS="$IFS" ; local IFS=$'\n' # Change divider for iterator one line below
    opts=$(for x in `udfs files ls ${lastDir}`; do echo ${lastDir}${x}/ ; done) # TODO: Implement "udfs files ls -F" to get rid of frontslash after files. This does currently throw "Error: /cats/foo/ is not a directory"
    IFS="$OLDIFS" # Reset divider to space, ' '
    COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
    [[ $COMPREPLY = */ ]] && compopt -o nospace -o filenames
    return 0
}

_udfs_multiaddr_complete()
{
    local lastDir=${word%/*}/
    # Special case
    if [[ ${word} == */"ipcidr"* ]] ; then # TODO: Broken, fix it.
        opts="1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32" # TODO: IPv6?
        COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
    # "Loop"
    elif [[ ${word} == /*/ ]] || [[ ${word} == /*/* ]] ; then
        if [[ ${word} == /*/*/*/*/*/ ]] ; then
            COMPREPLY=()
        elif [[ ${word} == /*/*/*/*/ ]] ; then
            word=${word##*/}
            opts="udfs/ "
            COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
        elif [[ ${word} == /*/*/*/ ]] ; then
            word=${word##*/}
            opts="4001/ "
            COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
        elif [[ ${word} == /*/*/ ]] ; then
            word=${word##*/}
            opts="udp/ tcp/ ipcidr/"
            COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
        elif [[ ${word} == /*/ ]] ; then
            COMPREPLY=() # TODO: This need to return something to NOT break the function. Maybe a "/" in the end as well due to -o filename option.
        fi
        COMPREPLY=${lastDir}${COMPREPLY}
    else # start case
        opts="/ip4/ /ip6/"
        COMPREPLY=( $(compgen -W "${opts}" -- ${word}) )
    fi
    [[ $COMPREPLY = */ ]] && compopt -o nospace -o filenames
    return 0
}

_udfs_pinned_complete()
{
    local OLDIFS="$IFS" ; local IFS=$'\n'
    local pinned=$(udfs pin ls)
    COMPREPLY=( $(compgen -W "${pinned}" -- ${word}) )
    IFS="$OLDIFS"
    if [[ ${#COMPREPLY[*]} -eq 1 ]]; then # Only one completion, remove pretty output
        COMPREPLY=( ${COMPREPLY[0]/ *//} ) #Remove ' ' and everything after
        [[ $COMPREPLY = */ ]] && compopt -o nospace  # Removing whitespace after output
    fi
}
_udfs_filesystem_complete()
{
    compopt -o default # Re-enable default file read
    COMPREPLY=()
}

_udfs()
{
    COMPREPLY=()
    compopt +o default # Disable default to not deny completion, see: http://stackoverflow.com/a/19062943/1216348

    local word="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"

    case "${COMP_CWORD}" in
        1)
            local opts="add bitswap block bootstrap cat commands config daemon dag dht \
                        diag dns file files get id init log ls mount name object pin ping pubsub \
                        refs repo resolve stats swarm tar update version"
            COMPREPLY=( $(compgen -W "${opts}" -- ${word}) );;
        2)
            local command="${COMP_WORDS[1]}"
            eval "_udfs_$command" 2> /dev/null ;;
        *)
            local command="${COMP_WORDS[1]}"
            local subcommand="${COMP_WORDS[2]}"
            eval "_udfs_${command}_${subcommand}" 2> /dev/null && return
            eval "_udfs_$command" 2> /dev/null ;;
    esac
}
complete -F _udfs udfs
