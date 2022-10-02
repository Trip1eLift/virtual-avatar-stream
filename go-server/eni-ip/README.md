Docs: https://stackoverflow.com/questions/49354116/how-do-i-retrieve-the-public-ip-for-a-fargate-task-using-the-cli

List ENI managed by amazon
```
aws ec2 describe-network-interfaces --filter Name=association.ip-owner-id,Values=amazon | grep -i PublicIp
```

Operations:
```
export TASK_ARN=$(aws ecs list-tasks --cluster "virtual-avatar-stream-dev-cluster" --query 'taskArns[0]' --output text)

export ENI=$(aws ecs describe-tasks --cluster "virtual-avatar-stream-dev-cluster" --task "${TASK_ARN}" --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' --output text)

aws ec2 describe-network-interfaces --network-interface-ids "${ENI}" --query 'NetworkInterfaces[0].Association.PublicIp' --output text
```