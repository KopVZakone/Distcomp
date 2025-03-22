package issue

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Khmelov/Distcomp/251003/Nasevich/restbasics/internal/storage/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrIssueNotFound     = fmt.Errorf("issue not found")
	ErrFailedToCreate    = fmt.Errorf("failed to create issue")
	ErrFailedToUpdate    = fmt.Errorf("failed to update issue")
	ErrFailedToDelete    = fmt.Errorf("failed to delete issue")
	ErrInvalidIssueData  = fmt.Errorf("invalid issue data")
	ErrInvalidForeignKey = fmt.Errorf("invalid foreign key passed")
)

type Issue interface {
	CreateIssue(ctx context.Context, is model.Issue) (model.Issue, error)
	GetIssues(ctx context.Context) ([]model.Issue, error)
	GetIssueByID(ctx context.Context, id int64) (model.Issue, error)
	UpdateIssueByID(ctx context.Context, is model.Issue) (model.Issue, error)
	DeleteIssueByID(ctx context.Context, id int64) error
}

type instance struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Issue {
	return &instance{
		db: db,
	}
}

func (i *instance) CreateIssue(ctx context.Context, is model.Issue) (model.Issue, error) {
	query := `INSERT INTO Issue (creatorid, title, content) 
	          VALUES ($1, $2, $3) RETURNING id, created`

	var id int64
	var created time.Time

	err := i.db.QueryRowContext(ctx, query, is.CreatorID, is.Title, is.Content).
		Scan(&id, &created)
	if err != nil {
		log.Println(err)

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return model.Issue{}, ErrIssueNotFound
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return model.Issue{}, ErrInvalidIssueData
		}

		return model.Issue{}, fmt.Errorf("failed to create issue: %w", err)
	}

	is.ID = id
	is.Created = created

	return is, nil
}

func (i *instance) GetIssues(ctx context.Context) ([]model.Issue, error) {
	var issues []model.Issue
	query := `SELECT * FROM Issue`

	err := i.db.SelectContext(ctx, &issues, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve issues: %w", err)
	}

	if len(issues) == 0 {
		return []model.Issue{}, nil
	}

	return issues, nil
}

func (i *instance) GetIssueByID(ctx context.Context, id int64) (model.Issue, error) {
	var issue model.Issue
	query := `SELECT id, creatorId, title, content, created, modified FROM Issue WHERE id = $1`

	err := i.db.GetContext(ctx, &issue, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return issue, ErrIssueNotFound
		}
		return issue, fmt.Errorf("failed to retrieve issue by ID: %w", err)
	}

	return issue, nil
}

func (i *instance) UpdateIssueByID(ctx context.Context, is model.Issue) (model.Issue, error) {
	query := `UPDATE Issue SET creatorId = $1, title = $2, content = $3, modified = $4 
	          WHERE id = $5 RETURNING id, creatorId, title, content, created, modified`
	var updatedIssue model.Issue

	err := i.db.QueryRowContext(ctx, query, is.CreatorID, is.Title, is.Content, is.Modified, is.ID).
		Scan(&updatedIssue.ID, &updatedIssue.CreatorID, &updatedIssue.Title, &updatedIssue.Content, &updatedIssue.Created, &updatedIssue.Modified)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("issue not found with id:", is.ID)
			return updatedIssue, ErrIssueNotFound
		}

		log.Println("error with query:", err)
		return updatedIssue, fmt.Errorf("failed to update issue: %w", err)
	}

	return updatedIssue, nil
}

func (i *instance) DeleteIssueByID(ctx context.Context, id int64) error {
	query := `DELETE FROM Issue WHERE id = $1`
	result, err := i.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("Error executing DELETE query:", err)
		return fmt.Errorf("failed to delete issue: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Println("No issue found with ID:", id)
		return ErrIssueNotFound
	}

	return nil
}
