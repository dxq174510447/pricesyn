package dynamic_plan

/*
动态规划
时间复杂度
1. 多项式时间复杂度，O(1),O(log(n)),O(n^a)
2. 非多项式时间复杂度 O(a^n)和O(n!)
np,非多项式问题，很难找出解，只能找出最优解
动态规划 先解决子问题，在逐渐解决大问题。（每一步找到最优解，达到结果是最优）
 */

type Product struct {
	Name string
	// 都按1的维度
	Weight int
}

type DynamicPlan struct {
	
}