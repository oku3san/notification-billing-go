package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "log"
    "os"
    "sync/atomic"
    "time"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/budgets"
    "github.com/slack-go/slack"
)

type Response struct {
    ActualSpend     string `json:"actual_spend"`
    ForecastedSpend string `json:"forecasted_spend"`
}

func handler(ctx context.Context) (Response, error) {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return Response{}, err
    }

    svc := budgets.NewFromConfig(cfg)

    input := &budgets.DescribeBudgetInput{
        AccountId:  aws.String(os.Getenv("AccountID")),
        BudgetName: aws.String(os.Getenv("BudgetName")),
    }

    result, err := svc.DescribeBudget(ctx, input)
    if err != nil {
        return Response{}, err
    }

    actualSpend := result.Budget.CalculatedSpend.ActualSpend.Amount
    forecastedSpend := result.Budget.CalculatedSpend.ForecastedSpend.Amount

    if err := sendNotification(*actualSpend, *forecastedSpend); err != nil {
        return Response{}, err
    }

    response := Response{
        ActualSpend:     *actualSpend,
        ForecastedSpend: *forecastedSpend,
    }

    return response, nil
}

// Slack APIの呼び出しが完了したかどうかを保持するためのフラグ
var slackAPICompleted int32

func sendNotification(actualSpend, forecastedSpend string) error {
    webhookURL := os.Getenv("WebhookURL")
    message := fmt.Sprintf("実績値: $%s\n予測値: $%s", actualSpend, forecastedSpend)

    payload := slack.WebhookMessage{
        Text: message,
    }

    // Slack APIの呼び出しを非同期的に実行
    go func() {
        err := slack.PostWebhook(webhookURL, &payload)
        if err != nil {
            log.Printf("error: %s", err)
        }

        // Slack APIの呼び出しが完了したことを通知
        atomic.StoreInt32(&slackAPICompleted, 1)
    }()

    // Slack APIの呼び出しが完了するまで待機
    for {
        if atomic.LoadInt32(&slackAPICompleted) == 1 {
            return nil
        }

        // 100ミリ秒待機
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    lambda.Start(handler)
}
