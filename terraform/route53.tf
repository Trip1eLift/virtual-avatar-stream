# TODO: Fix connection timeout from domain name

resource "aws_route53_record" "root-a" {
  zone_id = var.hosted_zone_id
  name    = var.domain_name
  type    = "A"
  # ttl  = 300 # TTL for all alias records is 60 seconds, you cannot change this, therefore ttl has to be omitted in alias records.
  alias {
    name    = aws_lb.main.dns_name
    zone_id = aws_lb.main.zone_id
    evaluate_target_health = false
  }
}