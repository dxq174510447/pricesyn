package dynamic_plan

import (
	"context"
	"fmt"
	"pricesyn/util"
	"strings"
)

/*
动态规划
时间复杂度
1. 多项式时间复杂度，O(1),O(log(n)),O(n^a)
2. 非多项式时间复杂度 O(a^n)和O(n!)
np,非多项式问题，很难找出解，只能找出最优解
动态规划 先解决子问题，在逐渐解决大问题。（每一步找到最优解，达到结果是最优）
都是从表格开始
典型案例：背包问题，最大子串问题
*/

type Thing struct {
	Name   string
	Size   int
	Weight int
}

type ColumnDef struct {
	Desc     string
	Capacity int
}

type CellDef struct {
	X int
	Y int
	// 当前cell所包含的权重
	Weight int
	// 当前cell所包含的thing大小
	Size    int
	Filling []*Thing
}

func (c *CellDef) Print() string {
	fillstr := util.StringUtil.ArrayJoin(c.Filling, func(index int) string {
		return c.Filling[index].Name
	})
	msg := fmt.Sprintf("%d,%d{weight %d size %d ,filling %s }", c.X, c.Y, c.Weight, c.Size, fillstr)
	return msg
}

func (c *CellDef) Fill(t *Thing) {
	c.Weight = c.Weight + t.Weight
	c.Size = c.Size + t.Size
	c.Filling = append(c.Filling, t)
}

type DynamicPlan struct {
	X     []*ColumnDef
	Y     []*Thing
	Cells [][]*CellDef
}

func (d *DynamicPlan) Print() {
	for _, row := range d.Cells {
		var msg []string
		for _, cell := range row {
			msg = append(msg, cell.Print())
		}
		fmt.Println(strings.Join(msg, "|"))
	}
}

func (d *DynamicPlan) newCellDef(x int, y int, weight int, size int, thing *Thing) *CellDef {
	cell := &CellDef{
		X:      x,
		Y:      y,
		Weight: weight,
		Size:   size,
	}
	if thing == nil {
		return cell
	}
	cell.Filling = append(cell.Filling, thing)
	return cell
}

func (d *DynamicPlan) copyCellDef(cell *CellDef) *CellDef {
	newCell := &CellDef{
		X:      cell.X,
		Y:      cell.Y,
		Size:   cell.Size,
		Weight: cell.Weight,
	}
	if len(cell.Filling) == 0 {
		return newCell
	}
	target := make([]*Thing, len(cell.Filling), len(cell.Filling))
	copy(target, cell.Filling)
	newCell.Filling = target
	return newCell
}

func (d *DynamicPlan) Run(ctx context.Context) {
	d.Cells = make([][]*CellDef, len(d.Y), len(d.Y))
	for i := 0; i < len(d.Y); i++ {
		thing := d.Y[i]
		d.RunRow(ctx, thing, i)
	}
}

func (d *DynamicPlan) RunRow(ctx context.Context, thing *Thing, row int) {

	d.Cells[row] = make([]*CellDef, len(d.X), len(d.X))
	for j := 0; j < len(d.X); j++ {
		col := d.X[j]
		if col.Capacity < thing.Size {
			if row == 0 {
				d.Cells[row][j] = d.newCellDef(j, row, 0, 0, nil)
			} else {
				d.Cells[row][j] = d.copyCellDef(d.Cells[row-1][j])
			}
		} else {
			if row == 0 {
				d.Cells[row][j] = d.newCellDef(j, row, thing.Weight, thing.Size, thing)
			} else {
				upperCell := d.Cells[row-1][j]

				leftCapacity := col.Capacity - thing.Size
				var lindex int = -1
				if leftCapacity > 0 {
					for lastIndex, lastCol := range d.X {
						if lastCol.Capacity >= leftCapacity {
							lindex = lastIndex
							break
						}
					}
				}

				var currentWeight int = 0
				if lindex == -1 {
					currentWeight = thing.Weight
				} else {
					currentWeight = thing.Weight + d.Cells[row-1][lindex].Weight
				}

				if upperCell.Weight > currentWeight {
					// 不替换
					d.Cells[row][j] = d.copyCellDef(d.Cells[row-1][j])
				} else {
					// 要替换
					if lindex == -1 {
						d.Cells[row][j] = d.newCellDef(j, row, thing.Weight, thing.Size, thing)
					} else {
						d.Cells[row][j] = d.newCellDef(j, row, thing.Weight, thing.Size, thing)
						for _, pre := range d.Cells[row-1][lindex].Filling {
							d.Cells[row][j].Fill(pre)
						}
					}
				}
			}
		}
	}
}

func (d *DynamicPlan) AddThing(t *Thing) *DynamicPlan {
	d.Y = append(d.Y, t)
	return d
}

func (d *DynamicPlan) AddColumn(col *ColumnDef) *DynamicPlan {
	d.X = append(d.X, col)
	return d
}

func NewDynamicPlan() *DynamicPlan {
	return &DynamicPlan{}
}
