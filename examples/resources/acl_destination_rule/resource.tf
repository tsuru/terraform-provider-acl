# first: create the instance acl via tsuru/tsuru provider
resource "tsuru_service_instance" "acl" {
  service_name = "acl"
  name         = "<< APP_NAME >>"
  owner        = "<< TEAM_NAME >>"
  plan         = ""
}

# second: bind previous created instance with tsuru app
resource "tsuru_service_instance_bind" "app-acl" {
  service_name     = tsuru_service_instance.acl.service_name
  service_instance = tsuru_service_instance.acl.name
  app              = "<< APP_NAME >>"
}



# first scenario, a app accessing another app
resource "acl_destination_rule" "test_app" {
  instance = tsuru_service_instance.acl.name

  app = "<< DESTINATION-APP >>"
}


# second scenario, a app accessing a tsuru reverse proxy instance
resource "acl_destination_rule" "test_app" {
  instance = tsuru_service_instance.acl.name

  rpaas {
    service_name = "<< MY-RPAAS-SERVICE >>"
    instance     = "<< MY-RPAAS-INSTANCE >>"
  }
}

# third scenario, a app accessing a external service via DNS
resource "acl_destination_rule" "test_app" {
  instance = tsuru_service_instance.acl.name

  dns = "example.org"

  port {
    number   = 80
    protocol = "TCP"
  }
}


# fourth scenario, a app accessing a external a network
resource "acl_destination_rule" "test_app" {
  instance = tsuru_service_instance.acl.name

  ip = "<< NETWORK CIDR >>"

  port {
    number   = 80
    protocol = "TCP"
  }

  port {
    number   = 443
    protocol = "TCP"
  }
}
