# for app
terraform import acl_destination_rule.resource_name "service::acl::instance::app::app-name"

# example
terraform import acl_destination_rule.my_acl_app "acl-rule::acl::my-acl::app::sample-app"

# for dns
terraform import acl_destination_rule.resource_name "service::acl::instance::dns::dns-name"

# example
terraform import acl_destination_rule.my_acl_dns "acl-rule::acl::my-acl::dns::example.com"

# for ip
terraform import acl_destination_rule.resource_name "service::acl::instance::ip::cidr-target"

# example
terraform import acl_destination_rule.my_acl_ip "acl-rule::acl::my-acl::ip::10.0.0.1/24"

# for rpaas
terraform import acl_destination_rule.resource_name "service::acl::instance::rpaas::service-instance::instance-rpaas"

# example for rpaasv2-be
terraform import acl_destination_rule.my_acl_rpaas "acl-rule::acl::my-acl::rpaas::rpaasv2-be::sample-app-rpaas"
