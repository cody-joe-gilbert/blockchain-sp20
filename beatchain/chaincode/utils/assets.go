package utils

type Transaction struct {
	/*
	Defines the standard format of the user's transaction info. Used to prevent having to
	replicate boilerplate code on each transaction function
	 */
	CalledFunction string
	CreatorOrg string
	CreatorCertIssuer string
	Args []string
}
