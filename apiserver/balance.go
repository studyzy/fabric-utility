package main

import (
	"errors"
	"math/rand"
)

type Balancer interface {
	DoBalance() (string, error)
}

type BalancerMgr struct {
	OrgBalancers map[string]Balancer
}

func (b *BalancerMgr) RegisterBalancer(orgName string, balancer Balancer) error {
	if b == nil {
		return errors.New("balancerMgr null. ")
	}
	b.OrgBalancers[orgName] = balancer
	return nil
}

// =============================================================
type RandomBalancer struct {
	OrgName string
	Peers   []string
}

func (rd *RandomBalancer) DoBalance() (string, error) {
	if rd == nil {
		return "", errors.New("randomBalancer null. ")
	}
	num := len(rd.Peers)
	if num == 0 {
		return "", errors.New("no peers in randomBalancer. ")
	}
	index := rand.Intn(num)
	peer := rd.Peers[index]
	return peer, nil
}

// =============================================================
type RoundRobinBalancer struct {
	OrgName string
	Peers   []string
	Current int
}

func (rr *RoundRobinBalancer) DoBalance() (string, error) {
	if rr == nil {
		return "", errors.New("roundRobinBalancer null. ")
	}
	num := len(rr.Peers)
	if num == 0 {
		return "", errors.New("no peers in roundRobinBalancer. ")
	}
	peer := rr.Peers[rr.Current]
	rr.Current++
	if rr.Current >= num {
		rr.Current = 0
	}
	return peer, nil
}
