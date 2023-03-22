package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/budgets"
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

    // budgetsサービスのクライアントを作成
    svc := budgets.NewFromConfig(cfg)

    // describeBudget APIを呼び出し、ActualSpendとForecastedSpendを取得
    input := &budgets.DescribeBudgetInput{
        AccountId:  aws.String(os.Getenv("AccountId")),
        BudgetName: aws.String(os.Getenv("BudgetName")),
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

    fmt.Println("aaaaaaaa")

    return response, nil
}

func main() {
    lambda.Start(handler)
}
