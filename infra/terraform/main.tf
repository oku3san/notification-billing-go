# IAM role for the Lambda function
resource "aws_iam_role" "lambda_role" {
  name = "lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Policy for the Lambda function to access AWS Budgets
resource "aws_iam_policy" "lambda_policy" {
  name_prefix = "lambda-policy-"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "budgets:DescribeBudgets",
          "budgets:ViewBudget"
        ]
        Resource = "*"
      }
    ]
  })
}

# Attach the policy to the Lambda role
resource "aws_iam_role_policy_attachment" "lambda_policy_attachment" {
  policy_arn = aws_iam_policy.lambda_policy.arn
  role       = aws_iam_role.lambda_role.name
}

# AWS Lambda function
resource "aws_lambda_function" "notification-billing-go" {
  filename      = "artifact/notification-billing-go.zip"
  function_name = "notification-billing-go"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2"
  publish       = true

  source_code_hash = filebase64sha256("artifact/notification-billing-go.zip")

  environment {
    variables = {
      "AccountID"  = var.AccountID
      "BudgetName" = var.BudgetName
      "WebhookURL" = var.WebhookURL
    }
  }
}

resource "aws_cloudwatch_event_rule" "notification-billing-go" {
  name                = "notification-billing-go"
  schedule_expression = "cron(30 15 * * ? *)"
}

resource "aws_cloudwatch_event_target" "notification-billing-go" {
  target_id = "notification-billing-go"
  arn       = aws_lambda_function.notification-billing-go.arn
  rule      = aws_cloudwatch_event_rule.notification-billing-go.name

}

resource "aws_lambda_permission" "notification-billing-go" {
  statement_id  = "notification-billing-go"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.notification-billing-go.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.notification-billing-go.arn
}
