package chart_demo

import (
	"github.com/r0vx/admin/presets"
	"github.com/r0vx/web"
	"github.com/r0vx/x/ui/shadcn"
	"github.com/r0vx/x/ui/unovis"
	h "github.com/r0vx/htmlgo"
	"gorm.io/gorm"
)

// NetworkGraphDemo 网络图演示
type NetworkGraphDemo struct{}

// configNetworkGraphDemo 配置网络图演示
func ConfigNetworkGraphDemo(pb *presets.Builder, db *gorm.DB) {
	b := pb.Model(&NetworkGraphDemo{}).Label("Network Graph Demo").URIName("network-graph-demo")

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		// 构建网络图数据（模拟用户提供的截图示例）
		graphData := unovis.NetworkGraphData{
			Nodes: []unovis.NetworkNodeDatum{
				// 中心节点：用户
				{
					ID:       "user-1",
					Label:    "jdoe@acme.com",
					SubLabel: "External User",
					Type:     "center",
					Color:    "hsl(217 91% 60%)", // 蓝色
					Badges: []unovis.NodeBadge{
						{Value: 3, Color: "hsl(217 91% 60%)"},   // 蓝色徽章
						{Value: 12, Color: "hsl(351 83% 82%)"},  // 粉红色徽章
						{Value: 150, Color: "hsl(48 96% 89%)"},  // 黄色徽章
					},
				},
				// 普通节点
				{
					ID:       "role-1",
					Label:    "AWSReservedSSO_Something",
					SubLabel: "Role",
					Color:    "hsl(217 91% 60%)", // 蓝色
					Badges: []unovis.NodeBadge{
						{Value: 2},
					},
				},
				{
					ID:       "instance-1",
					Label:    "i-0a1b...f6g7h8",
					SubLabel: "EC2 Instance",
					Color:    "hsl(351 83% 82%)", // 粉红色
					Badges: []unovis.NodeBadge{
						{Value: 1},
					},
				},
				{
					ID:       "instance-2",
					Label:    "i-1a1b...f6g7h8",
					SubLabel: "EC2 Instance",
					Color:    "hsl(351 83% 82%)", // 粉红色
				},
				{
					ID:       "file-1",
					Label:    "my-file",
					SubLabel: "File",
					Color:    "hsl(48 96% 89%)", // 黄色
				},
				{
					ID:       "secret-1",
					Label:    "tests-...ansfer",
					SubLabel: "Secret",
					Color:    "hsl(142 76% 73%)", // 绿色
				},
			},
			Links: []unovis.NetworkLinkDatum{
				{Source: "user-1", Target: "role-1"},
				{Source: "user-1", Target: "instance-1"},
				{Source: "user-1", Target: "instance-2"},
				{Source: "role-1", Target: "file-1"},
				{Source: "file-1", Target: "secret-1"},
				{Source: "instance-1", Target: "file-1"},
			},
		}

		// 配置节点
		graphConfig := unovis.NetworkGraphConfig{
			"user-1": unovis.NetworkNodeConfig{
				Label:    "jdoe@acme.com",
				SubLabel: "External User",
				Type:     "center",
				Color:    "hsl(217 91% 60%)",
				Size:     80,
				Tooltip: &unovis.NodeTooltip{
					Title: "User Details",
					Fields: []unovis.TooltipField{
						{Label: "Email", Value: "jdoe@acme.com"},
						{Label: "Type", Value: "External User"},
						{Label: "Status", Value: "Active"},
						{Label: "Last Login", Value: "2 hours ago"},
					},
				},
			},
			"role-1": unovis.NetworkNodeConfig{
				Label:    "AWSReservedSSO_Something",
				SubLabel: "Role",
				Tooltip: &unovis.NodeTooltip{
					Title: "Role Information",
					Fields: []unovis.TooltipField{
						{Label: "Name", Value: "AWSReservedSSO_Something"},
						{Label: "Type", Value: "IAM Role"},
						{Label: "Permissions", Value: "Read, Write"},
					},
				},
			},
			"instance-1": unovis.NetworkNodeConfig{
				Label:    "i-0a1b...f6g7h8",
				SubLabel: "EC2 Instance",
				Tooltip: &unovis.NodeTooltip{
					Title: "EC2 Instance",
					Fields: []unovis.TooltipField{
						{Label: "Instance ID", Value: "i-0a1b...f6g7h8"},
						{Label: "Type", Value: "t2.micro"},
						{Label: "Status", Value: "Running"},
						{Label: "Region", Value: "us-east-1"},
					},
				},
			},
			"instance-2": unovis.NetworkNodeConfig{
				Label:    "i-1a1b...f6g7h8",
				SubLabel: "EC2 Instance",
				Tooltip: &unovis.NodeTooltip{
					Title: "EC2 Instance",
					Fields: []unovis.TooltipField{
						{Label: "Instance ID", Value: "i-1a1b...f6g7h8"},
						{Label: "Type", Value: "t2.small"},
						{Label: "Status", Value: "Running"},
						{Label: "Region", Value: "us-west-2"},
					},
				},
			},
			"file-1": unovis.NetworkNodeConfig{
				Label:    "my-file",
				SubLabel: "File",
				Tooltip: &unovis.NodeTooltip{
					Title: "File Details",
					Fields: []unovis.TooltipField{
						{Label: "Name", Value: "my-file.txt"},
						{Label: "Size", Value: "2.5 MB"},
						{Label: "Modified", Value: "1 day ago"},
					},
				},
			},
			"secret-1": unovis.NetworkNodeConfig{
				Label:    "tests-...ansfer",
				SubLabel: "Secret",
				Tooltip: &unovis.NodeTooltip{
					Title: "Secret Information",
					Fields: []unovis.TooltipField{
						{Label: "Name", Value: "tests-transfer-secret"},
						{Label: "Type", Value: "API Key"},
						{Label: "Expires", Value: "30 days"},
					},
				},
			},
		}

		// 创建网络图
		graph := unovis.NetworkGraph().
			Data(graphData).
			LinkDistance(200).
			LinkStrength(0.5).
			Charge(-800).
			NodeSize(40).
			ShowLabels(true).
			Class("w-full h-[600px]")

		// 设置节点配置
		for id, cfg := range graphConfig {
			graph.NodeConfig(id, cfg)
		}

		body := h.Div(
			h.H1("Network Graph Demo").Class("text-2xl font-bold mb-4"),
			h.P(h.Text("展示用户与资源的关系网络图，支持自定义节点样式和工具提示。")).Class("text-muted-foreground mb-6"),

			// 网络图
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("用户资源关系图")),
					shadcn.CardDescription(h.Text("Force Layout 力导向布局，鼠标悬停查看详细信息")),
				),
				shadcn.CardContent(
					graph,
				),
			).Class("mb-6"),

			// 说明文档
			shadcn.Card(
				shadcn.CardHeader(
					shadcn.CardTitle(h.Text("功能说明")),
				),
				shadcn.CardContent(
					h.Ul(
						h.Li(h.Text("支持 Force Layout 力导向布局，自动计算节点位置")),
						h.Li(h.Text("支持自定义节点颜色、大小和标签")),
						h.Li(h.Text("支持工具提示，鼠标悬停显示详细信息")),
						h.Li(h.Text("支持节点之间的连接关系")),
						h.Li(h.Text("中心节点可以配置为更大尺寸，突出显示")),
					).Class("list-disc list-inside space-y-2 text-sm text-muted-foreground"),
				),
			),
		).Class("container mx-auto p-4")

		r.Body = body
		r.PageTitle = "Network Graph Demo"

		return
	})
}
