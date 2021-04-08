package server

// APIResponse - API 응답 구조
type APIResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

// KubeConfigRequest - Kube Config 요청 구조
type KubeConfigRequest struct {
	Name       string `json:"name"`
	KubeConfig string `json:"kubeConfig"`
}

// KubeTokenRequest - Kube Token 요청 구조
type KubeTokenRequest struct {
	Name      string `json:"name"`
	APIServer string `json:"apiServer"`
	Token     string `json:"token"`
}
