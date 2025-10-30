package auth

import (
	"ticketing-infra/rpc-client/model"
	"ticketing-infra/rpc-client/response"

	"github.com/golang-jwt/jwt/v4"
)

func RefreshTokenHandler(refreshToken string) (model.Tokens, error) {
	// 验证刷新令牌
	token, err := ValidateRefreshToken(refreshToken)
	if err != nil || !token.Valid { // 刷新令牌无效
		return model.Tokens{}, response.InvalidRefreshToken
	}

	// 生成新的访问令牌和刷新令牌
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := int(claims["user_id"].(float64))
		newAccessToken, newRefreshToken, err := GenerateTokens(userID)
		if err != nil {
			return model.Tokens{}, err
		}

		// 返回新的访问令牌和刷新令牌
		return model.Tokens{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
	} else {
		return model.Tokens{}, response.InvalidClaims
	}
}
