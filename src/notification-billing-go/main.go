package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/budgets"
    "github.com/slack-go/slack"
    "log"
    "os"
)

type Response struct {
    ActualSpend     string `json:"actual_spend"`
    ForecastedSpend string `json:"forecasted_spend"`
}

func handler(ctx context.Context) (Response, error) {
    // AWS SDK for Go v2の設定を読み込む
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return Response{}, fmt.Errorf("failed to load SDK configuration, %v", err)
    }

    accountId := os.Getenv("AccountID")
    budgetName := os.Getenv("BudgetName")

    // budgetsサービスのクライアントを作成
    svc := budgets.NewFromConfig(cfg)

    // describeBudget APIを呼び出し、ActualSpendとForecastedSpendを取得
    input := &budgets.DescribeBudgetInput{
        AccountId:  &accountId,
        BudgetName: &budgetName,
    }

    result, err := svc.DescribeBudget(ctx, input)
    if err != nil {
        return Response{}, fmt.Errorf("failed to describe budget, %v", err)
    }

    actualSpend := result.Budget.CalculatedSpend.ActualSpend.Amount
    forecastedSpend := result.Budget.CalculatedSpend.ForecastedSpend.Amount

    response := Response{
        ActualSpend:     *actualSpend,
        ForecastedSpend: *forecastedSpend,
    }

    sendNotification(response)

    return response, nil
}

func sendNotification(response Response) error {
    webhookURL := os.Getenv("WebhookURL")
    message := fmt.Sprintf("実績値: $%s\n予測値: $%s", response.ActualSpend, response.ForecastedSpend)

    payload := slack.WebhookMessage{
        Text: message,
    }

    err := slack.PostWebhook(webhookURL, &payload)
    if err != nil {
        log.Fatalf("error: %s", err)
    }

    return nil
}

func main() {
    lambda.Start(handler)
}
