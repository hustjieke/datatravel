./datatravel --threads=128 --checksum=true --max-allowed-packet-MB=16 --fk-check=false --debug=false --to-flavor=radondb --set-global-read-lock=true --from=172.31.1.2:3306 --from-user=superadmin --from-password=superpasswd --to=172.31.8.30:3306 --to-user=radon --to-password=Sun55_kongg --table-db=ydcpbxt --tables=t_lpr_result,t_engine_result > result_datatravel.log 2>&1 &
