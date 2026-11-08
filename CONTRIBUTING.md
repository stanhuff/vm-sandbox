# Contributing

## Test coverage

Target at least 90% aggregate statement coverage for maintained production code. Set up CI coverage reporting with the first implementation PR and check the target on every PR. Use unit and integration tests to cover behavior, including denial cases, quota boundaries, concurrency, and failure recovery as those features land.

If the target cannot reasonably be met, document the specific gap and reason in the PR. Keep exclusions narrow and visible in CI. Make any temporary lower threshold explicit and raise it as coverage improves until it reaches 90%.
