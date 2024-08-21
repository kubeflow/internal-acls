"""A python program to check the github org file.

The purpose of this file is to guard against bad changes (e.g. adding someone)
as an admin who shouldn't be an admin.
"""
import fire
import logging
import yaml

class CheckConfig(object):
  def check_config(seld, config):
    """Check that the config is valid

    Args:
      config: Path to YAML file
    """
    with open(config) as hf:
      org = yaml.load(hf, Loader=yaml.Loader)

    admins = org.get("orgs").get("kubeflow").get("admins")

    # There should be at least some admins
    if not admins:
      error = "config {0} is not valid; missing orgs.kubeflow.admins".format(
        config)

      logging.error(error)
      raise ValueError(error)

    # TODO(jlewi): We should load this in via config map
    # Check that each admin is in a whitelist set of admins.
    allowed_admins = ["andreyvelich", "caniszczyk", "chensun", "googlebot",
                      "google-oss-robot", "james-jwu", "jbottum", "jlewi", "johnugeorge",
                      "k8s-ci-robot", "theadactyl", "krook", "thelinuxfoundation",
                      "terrytangyuan", "zijianjoy"]

    for a in admins:
      if not a in allowed_admins:
        error = ("{0} is not in the allowed set of admins. "
                 "Allowed admins is {1}").format(a, ", ".join(allowed_admins))
        logging.error(error)
        raise ValueError(error)

    logging.info("config is valid")

if __name__ == "__main__":
  logging.basicConfig(level=logging.INFO,
                    format=('%(levelname)s|%(asctime)s'
                            '|%(message)s|%(pathname)s|%(lineno)d|'),
                    datefmt='%Y-%m-%dT%H:%M:%S',
                    )

  fire.Fire(CheckConfig)
