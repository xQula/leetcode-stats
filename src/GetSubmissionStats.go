package src

import (
	"context"
	"math"

	"github.com/machinebox/graphql"
)

type SubmissionData struct {
	AllQuestionsCount []struct {
		Difficulty string `json:"difficulty"`
		Count      int    `json:"count"`
	} `json:"allQuestionsCount"`
	MatchedUser struct {
		ProblemsSolvedBeatsStats []struct {
			Difficulty string  `json:"difficulty"`
			Percentage float64 `json:"percentage"`
		} `json:"problemsSolvedBeatsStats"`
		SubmitStatsGlobal struct {
			AcSubmissionNum []struct {
				Difficulty string  `json:"difficulty"`
				Count      int     `json:"count"`
				Percentage float64 `json:"percentage"`
			} `json:"acSubmissionNum"`
		} `json:"submitStatsGlobal"`
	} `json:"matchedUser"`
}

func GetSubmissionStats(username string) (SubmissionData, error) {
	client := graphql.NewClient("https://leetcode.com/graphql")

	req := graphql.NewRequest(`
    query userProblemsSolved($username: String!) {
		allQuestionsCount {
		  difficulty
		  count
		}
		matchedUser(username: $username) {
		  problemsSolvedBeatsStats {
			difficulty
			percentage
		  }
		  submitStatsGlobal {
			acSubmissionNum {
			  difficulty
			  count
			}
		  }
		}
	  }
  `)
	req.Var("username", username)

	ctx := context.Background()

	var respData SubmissionData
	if err := client.Run(ctx, req, &respData); err != nil {
		return respData, err
	}

	for i := 0; i < len(respData.MatchedUser.ProblemsSolvedBeatsStats); i++ {
		respData.MatchedUser.ProblemsSolvedBeatsStats[i].Percentage = math.Round(respData.MatchedUser.ProblemsSolvedBeatsStats[i].Percentage*10) / 10
	}

	for i := 0; i < len(respData.MatchedUser.SubmitStatsGlobal.AcSubmissionNum); i++ {
		percentage := float64(respData.MatchedUser.SubmitStatsGlobal.AcSubmissionNum[i].Count) / float64(respData.AllQuestionsCount[i].Count)
		respData.MatchedUser.SubmitStatsGlobal.AcSubmissionNum[i].Percentage = math.Round(percentage*1000000) / 10000
	}

	return respData, nil
}