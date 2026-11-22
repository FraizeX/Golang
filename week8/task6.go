func buildReport(ctx context.Context) (*Report, error) {
    select {
    case <-time.After(20 * time.Millisecond):
        return &Report{}, nil

    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func handleFixed(timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(
        context.Background(), timeout,
    )
    defer cancel()

    ch := make(chan *Report, 1)

    go func() {
        report, err := buildReport(ctx)
        if err != nil {
            return
        }

        select {
        case ch <- report:
        case <-ctx.Done():
        }
    }()

    select {
    case <-ch:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}