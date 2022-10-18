# # ALB might be required to pass the provisioning state for ECS task?
# # https://stackoverflow.com/questions/63123466/all-tasks-on-a-ecs-service-stuck-in-provisioning-state

# resource "aws_lb" "main" {
#     name               = "${var.name}-${var.environment}-alb"
#     internal           = false
#     load_balancer_type = "application"
#     security_groups    = [ aws_security_group.ecs_tasks.id ]
#     subnets            = flatten(aws_subnet.public.*.id)

#     enable_deletion_protection = false
#     tags = var.common_tags
# }

# resource "aws_alb_target_group" "main" {
#     name        = "${var.name}-${var.environment}-tg"
#     port        = 5001
#     protocol    = "HTTP"
#     vpc_id      = aws_vpc.main.id
#     target_type = "ip"

#     health_check {
#         healthy_threshold   = "3"
#         interval            = "30"
#         protocol            = "HTTP"
#         matcher             = "200"
#         timeout             = "3"
#         path                = "/health"
#         unhealthy_threshold = "2"
#     }
#     tags = var.common_tags
# }

# # resource "aws_alb_listener" "http" {
# #   load_balancer_arn = aws_lb.main.id
# #   port              = 80
# #   protocol          = "HTTP"

# #   default_action {
# #    type = "redirect"

# #    redirect {
# #      port        = 443
# #      protocol    = "HTTPS"
# #      status_code = "HTTP_301"
# #    }
# #   }
# # }

# resource "aws_alb_listener" "https" {
#     load_balancer_arn = aws_lb.main.id
#     port              = 5001
#     protocol          = "HTTP"

# #ssl_policy        = "ELBSecurityPolicy-2016-08"
# #certificate_arn   = var.alb_tls_cert_arn

#     default_action {
#         target_group_arn = aws_alb_target_group.main.id
#         type             = "forward"
#     }
# }