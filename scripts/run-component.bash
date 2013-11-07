#!/dev/null

_identifier="${1:-000000006b1da14b35b70da1be305260915fda32}"
_fqdn="${mosaic_node_fqdn:-mosaic-0.loopback.vnet}"
_ip="${mosaic_node_ip:-127.0.155.0}"

if test -n "${mosaic_component_temporary:-}" ; then
	_tmp="${mosaic_component_temporary:-}"
elif test -n "${mosaic_temporary:-}" ; then
	_tmp="${mosaic_temporary}/components/${_identifier}"
else
	_tmp="/tmp/mosaic/components/${_identifier}"
fi

_run_bin="${_applications_elf}/mysqlb.elf"
_run_env=(
		mosaic_component_identifier="${_identifier}"
		mosaic_component_temporary="${_tmp}"
		mosaic_node_fqdn="${_fqdn}"
		mosaic_node_ip="${_ip}"
		transcript_level=information
)

if test "${_identifier}" == 00000000adf010553c8b9b9a8b3f78547f6b54c5 ; then
	if ! test "${#}" -ge 2 ; then
		echo "[ee] invalid arguments; aborting!" >&2
		exit 1
	fi
	_run_args=(
			component-init
			"${@:2}"
	)
elif test "${_identifier}" != 000000006b1da14b35b70da1be305260915fda32 ; then
	if ! test "${#}" -eq 1 ; then
		echo "[ee] invalid arguments; aborting!" >&2
		exit 1
	fi
	_run_args=(
			component "${_identifier}"
	)
else
	if ! test "${#}" -eq 0 ; then
		echo "[ee] invalid arguments; aborting!" >&2
		exit 1
	fi
	_run_args=(
			standalone
	)
fi

mkdir -p -- "${_tmp}"
cd -- "${_tmp}"

exec env "${_run_env[@]}" "${_run_bin}" "${_run_args[@]}"

exit 1
