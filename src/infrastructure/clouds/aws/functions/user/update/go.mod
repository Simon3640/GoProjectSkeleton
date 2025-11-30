module gormgoskeleton/functions/aws/user/update

go 1.25

require (
	github.com/aws/aws-lambda-go v1.47.0
	gormgoskeleton v0.0.0
	gormgoskeleton/aws v0.0.0
)

replace gormgoskeleton => ../../../../../../..

replace gormgoskeleton/aws => ../../..
