package service

import (
	"testing"

	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/model"
	"github.com/campusbooks/campusbooks/internal/repository"
)

func TestUserServiceProfileAndStats(t *testing.T) {
	env := newTestEnv(t)
	userRepo := repository.NewUserRepository(env.db)
	evalRepo := repository.NewEvaluationRepository(env.db)
	svc := NewUserService(userRepo, evalRepo, env.logger)

	u := &model.User{StudentNo: "U001", Email: "u@c.local", PasswordHash: "h", Role: constants.RoleStudent}
	if err := userRepo.Create(u); err != nil {
		t.Fatalf("seed: %v", err)
	}

	updated, err := svc.UpdateProfile(u.ID, dto.UpdateProfileRequest{
		Name: "张三", Department: "计算机学院", Campus: "东校区", Contact: "13800000000",
	})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if updated.Name != "张三" || updated.Department != "计算机学院" {
		t.Errorf("profile not updated: %+v", updated)
	}

	stats, err := svc.GetStats(u.ID)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.GoodRate != "100.0%" {
		t.Errorf("good_rate=%q, want 100.0%%", stats.GoodRate)
	}

	// add enough bad evaluations to flip the risk flag
	for i := 0; i < 3; i++ {
		e := &model.Evaluation{FromUserID: u.ID, ToUserID: u.ID, BookID: uint(i + 1), Type: constants.EvaluationBad}
		if err := env.db.Create(e).Error; err != nil {
			t.Fatalf("seed eval: %v", err)
		}
	}
	stats2, err := svc.GetStats(u.ID)
	if err != nil {
		t.Fatalf("stats2: %v", err)
	}
	if !stats2.RiskFlagged {
		t.Errorf("risk_flagged should be true with 100%% bad eval")
	}

	// list users
	list, err := svc.ListUsers("", 1, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("list total=%d, want 1", list.Total)
	}
}
