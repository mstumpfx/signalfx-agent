import string
from functools import partial as p

import pytest

from helpers.assertions import has_datapoint_with_dim, has_log_message, tcp_socket_open
from helpers.util import container_ip, run_service, run_agent, wait_for

pytestmark = [pytest.mark.collectd, pytest.mark.etcd, pytest.mark.monitor_with_endpoints]

MONITOR_CONFIG = string.Template(
    """
monitors:
- type: collectd/solr
  host: $host
  port: 8983
"""
)


def test_solr_monitor():
    with run_service("solr") as solr_container:
        host = container_ip(solr_container)
        config = MONITOR_CONFIG.substitute(host=host)
        assert wait_for(p(tcp_socket_open, host, 8983), 60), "service not listening on port"

        with run_agent(config) as [backend, get_output, _]:
            assert wait_for(p(has_datapoint_with_dim, backend, "plugin", "solr")), "Didn't get solr datapoints"
            assert not has_log_message(get_output().lower(), "error"), "error found in agent output!"
