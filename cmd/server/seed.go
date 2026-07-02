package main

import (
	"context"
	"log"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/domain"
	"git.epam.com/go-language-global-mentoring-program/internal/grpc/common/utils"
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	"github.com/brianvoe/gofakeit/v6"
)

func seedDatabase(store repo.Store) {
	ctx := context.Background()

	log.Println("start seeding our db")

	users := generateUsers(10)
	for _, user := range users {
		if err := store.CreateUser(ctx, user); err != nil {
			log.Printf("failed to seed our user: %v", err)
		}
	}
	log.Printf("created %d users:", len(users))

	projects := generateProjects(5)
	for _, project := range projects {
		if err := store.CreateProject(ctx, project); err != nil {
			log.Printf("failed to seed our project: %v", err)
		}
	}
	log.Printf("created %d projects", len(projects))

	issues := generateIssues(20, users, projects)
	for _, issue := range issues {
		if err := store.CreateIssue(ctx, issue); err != nil {
			log.Printf("failed to seed our issue: %v", err)
		}
	}
	log.Printf("created %d issues", len(issues))

	log.Println("seeding completed")
}

func generateUsers(count int) []*domain.User {
	var users []*domain.User

	for i := 0; i < count; i++ {
		user := &domain.User{
			ID:           utils.GenerateULID(),
			FirstName:    gofakeit.FirstName(),
			LastName:     gofakeit.LastName(),
			EmailAddress: gofakeit.Email(),
		}
		users = append(users, user)
	}

	return users
}

func generateProjects(count int) []*domain.Project {
	var projects []*domain.Project

	projectTemplates := []string{
		"EPAM Learn",
		"EPAM Campus",
		"Telescope",
		"EPAM DIAL",
		"GO BY EXAMPLE",
		"EFFECTIVE GO",
		"OUTLOOK",
		"GIT LAB",
		"Google",
		"UBER",
	}

	for i := 0; i < count; i++ {
		project := &domain.Project{
			ID:          utils.GenerateULID(),
			Name:        projectTemplates[i%len(projectTemplates)],
			Description: gofakeit.Sentence(8) + " " + gofakeit.Sentence(6),
		}
		projects = append(projects, project)
	}

	return projects
}

func generateIssues(count int, users []*domain.User, projects []*domain.Project) []*domain.Issue {
	var issues []*domain.Issue

	statuses := []domain.IssueStatus{
		domain.IssueStatusNew,
		domain.IssueStatusAssigned,
		domain.IssueStatusInProgress,
		domain.IssueStatusResolved,
		domain.IssueStatusClosed,
	}

	types := []domain.IssueType{
		domain.IssueTypeBug,
		domain.IssueTypeFeature,
		domain.IssueTypePerformance,
		domain.IssueTypeCosmetic,
	}

	priorities := []domain.Priority{
		domain.PriorityCritical,
		domain.PriorityMajor,
		domain.PriorityImportant,
		domain.PriorityMinor,
	}

	resolutions := []domain.IssueResolution{
		domain.ResolutionFixed,
		domain.ResolutionInvalid,
		domain.ResolutionWontFix,
		domain.ResolutionWorksForMe,
	}

	bugSummaries := []string{
		"Application crashes on unexpected input",
		"Data is not saved after user action",
		"Search returns incorrect or empty results",
		"Performance degrades under high load",
		"Authentication fails intermittently",
		"UI layout breaks on mobile devices",
		"API returns inconsistent responses",
		"Background jobs fail without error logs",
		"File upload fails for large files",
		"Notifications are delayed or not delivered",
	}

	featureSummaries := []string{
		"Add dark mode support",
		"Implement real-time notifications",
		"Add data export functionality",
		"Introduce user dashboard analytics",
		"Add advanced search capabilities",
		"Implement role-based access control",
		"Add mobile responsive UI improvements",
		"Introduce AI-assisted suggestions",
		"Add integration with external services",
		"Improve system performance and scalability",
	}

	for i := 0; i < count; i++ {
		status := statuses[gofakeit.Number(0, len(statuses)-1)]
		issueType := types[gofakeit.Number(0, len(types)-1)]
		priority := priorities[gofakeit.Number(0, len(priorities)-1)]

		project := projects[gofakeit.Number(0, len(projects)-1)]

		var summary string
		if issueType == domain.IssueTypeBug {
			summary = bugSummaries[gofakeit.Number(0, len(bugSummaries)-1)]
		} else {
			summary = featureSummaries[gofakeit.Number(0, len(featureSummaries)-1)]
		}

		var assigneeID string
		if len(users) > 0 {
			assignee := users[gofakeit.Number(0, len(users)-1)]
			assigneeID = assignee.ID
		}

		resolution := resolutions[gofakeit.Number(0, len(resolutions)-1)]
		if status == domain.IssueStatusResolved || status == domain.IssueStatusClosed {
			resolution = resolutions[gofakeit.Number(0, len(resolutions)-1)]
		}

		createDate := gofakeit.DateRange(time.Now().AddDate(0, -3, 0), time.Now().AddDate(0, 0, -1))
		modifyDate := gofakeit.DateRange(createDate, time.Now())

		issue := &domain.Issue{
			ID:          utils.GenerateULID(),
			CreateDate:  createDate,
			ModifyDate:  modifyDate,
			Summary:     summary,
			Description: gofakeit.Paragraph(2, 4, 8, " "),
			Status:      status,
			Resolution:  resolution,
			Type:        issueType,
			Priority:    priority,
			ProjectID:   project.ID,
			AssigneeID:  assigneeID,
		}

		issues = append(issues, issue)
	}

	return issues
}
