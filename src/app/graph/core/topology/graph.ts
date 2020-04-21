import * as d3Select			from "d3-selection";
import * as d3Force				from "d3-force";
import * as d3Drag				from "d3-drag";
import {GraphBase}				from "../graph.base";
import {LegendModel, Toolbar}	from "../toolbar";
import {ConfigModel}			from "../../model/models";
import {Topology as model}		from "../../model/graph.model"
import "./graph.css";

/**
 * Topology 그래프 랜더러
 * 
 * 참고
 * 		- 속성확인		: https://bl.ocks.org/steveharoz/8c3e2524079a8c440df60c1ab72b5d03
 * 		- Basic			: https://bl.ocks.org/mbostock/ad70335eeef6d167bc36fd3c04378048
 * 						: https://www.d3indepth.com/force-layout/
 * 		- Drag & Move	: https://bl.ocks.org/puzzler10/4438752bb93f45dc5ad5214efaa12e4a
 * 		- Highlight 	: https://medium.com/ninjaconcept/interactive-dynamic-force-directed-graphs-with-d3-da720c6d7811
 * 
 *  TODO
 * 		그룹처리
 * 			- [Force in Box](https://observablehq.com/@john-guerra/force-in-a-box)
 * 			- [Grouping nodes in a Force-Directed Graph](https://bl.ocks.org/XavierGimenez/a8e8c5e9aed71ba96bd52332682c0399)
 * 		기타
 * 			- https://stackoverflow.com/questions/47544041/how-to-visualize-groups-of-nodes-in-a-d3-force-directed-graph-layout/48001854
 */
export class TopologyGraph extends GraphBase {


	/**
	 * (abstract) 랜더링
	 * @param data 토플로지를 위한 k8s 데이터 (model.K8s)
	 */
	public populate(conf:ConfigModel.Config, svgEl:d3Select.Selection<SVGSVGElement, any, SVGElement, any>, bounds:DOMRect, outlineEl:d3Select.Selection<SVGGElement,any,SVGElement,any>) {

		// svg > defs
		if(svgEl.select("defs").size() == 0) svgEl.append("defs").call(TopologyGraph.renderDefs, conf);


		let data:model.Topology = conf.data;

		// 라인 추가
		let linksEl:d3Select.Selection<SVGElement,any,SVGElement,any>;
		if(data.links) {
			linksEl = outlineEl.selectAll("line.link")
				.data(data.links).enter()
					.append("line").attr("class",(d:model.Link) => {
						return (d.hidden)? `${d.kind} hidden`:`${d.kind} link`;
					})
		}

		// 노드 추가
		if(!data.nodes) return;
		let nodes:d3Select.Selection<SVGElement,any,SVGElement,any> = outlineEl.selectAll("g.node")
			.data(data.nodes).enter()
				.append("g")
				.attr("class","node")
				// .attr("id", d=>d.id)
				.call(TopologyGraph.renderNode);

		// 좌표 Simulation
		let forceLink:d3Force.ForceLink<d3Force.SimulationNodeDatum, d3Force.SimulationLinkDatum<d3Force.SimulationNodeDatum>>;
		forceLink = d3Force.forceLink()
			.id( (d:any) => d.id)
			.links(data.links)
 			.strength(1);

		let nodeSimulation = d3Force.forceSimulation(data.nodes)
			.velocityDecay(0.3)	// 그룹별 간격 (1~0 클수록 붙는다) 최적을 찾아야할 듯 (0.25, 0.35) 클수록 좁아진다.
			.force("gravity", d3Force.forceManyBody().strength(0.9))
			.force("charge", d3Force.forceManyBody().strength(-600))
			.force("theta", d3Force.forceManyBody().theta(0.01))
			.force("link", forceLink)
			.force('collision',  d3Force.forceCollide().radius(conf.topology.collision.radius))
			// .force('collide',  d3Force.forceCollide( (d:any) => { return d.kind=="cluster" ? 0 : 60}))
			.alpha(1)
			.alphaDecay(conf.topology.simulation.alphaDecay)	// ~0.0228 시뮬레이션 decay - 클수록 빠르지만 배치가 완벽하지 않음 (default:0.06)
			.force("center", d3Force.forceCenter(bounds.width/2, bounds.height/2));


		// tick
		nodeSimulation.on("tick", () => {
			nodes.attr("transform", d=> { return `translate(${d.x},${d.y})`; });

			linksEl.attr("x1", d=>d.source.x)
				.attr("y1", d=>d.source.y)
				.attr("x2", d=>d.target.x)
				.attr("y2", d=>d.target.y);
		});

		// onEnd 이벤트
		if (conf.topology.simulation.onEnd && typeof conf.topology.simulation.onEnd == "function") nodeSimulation.on("end", conf.topology.simulation.onEnd);
		


		// 범례
		const legends:Array<LegendModel> = [
			{
				header: "Element",
				rows: [
					{ label: "Cluster", ico: '<g><use class="ico" href="#ac_ic_node_cluster" width="20" height="20"></use></g>' },
					{ label: "Node", ico: '<g><use class="ico" href="#ac_ic_node_node" width="20" height="20"></use></g>' },
					{ label: "Pod", ico: '<g><use class="ico" href="#ac_ic_node_pod" width="20" height="20"></use></g>' },
					{ label: "Container", ico: '<g><use class="ico" href="#ac_ic_node_container" width="20" height="20"></use></g>' }
				]
			}
		]
		this.svg().call(Toolbar.render, this, legends);

	}

	/**
	 * defs 정의
	 * 
	 * @param defsEl def 엘리먼트
	 */
	private static renderDefs(defsEl:d3.Selection<SVGDefsElement, any, SVGElement, any>) {

		// https://material.io/resources/icons/

		let g = defsEl.append("symbol").attr("id", "ac_ic_node_node").attr("viewBox", "0 0 24 24").append("g")
			g.append("rect").attr("x","4").attr("y","2").attr("width","16").attr("height","20")
			g.append("path").attr("d", "M18,4v16H6V4H18 M18,2H6C4.9,2,4,2.9,4,4v16c0,1.1,0.9,2,2,2h12c1.1,0,2-0.9,2-2V4C20,2.9,19.1,2,18,2L18,2z M7,19h10v-6H7 V19z M10,10h4v1h3V5H7v6h3V10z")

		g = defsEl.append("symbol").attr("id", "ac_ic_node_pod").attr("viewBox", "0 0 24 24").append("g")
			g.append("rect").attr("x","4").attr("y","2").attr("width","16").attr("height","20")
			g.append("path").attr("d", "M19 3H5c-1.1 0-2 .9-2 2v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM5 10h3.13c.21.78.67 1.47 1.27 2H5v-2zm14 2h-4.4c.6-.53 1.06-1.22 1.27-2H19v2zm0-4h-5v1c0 1.07-.93 2-2 2s-2-.93-2-2V8H5V5h14v3zm-2 7h-3v1c0 .47-.19.9-.48 1.25-.37.45-.92.75-1.52.75s-1.15-.3-1.52-.75c-.29-.35-.48-.78-.48-1.25v-1H3v4c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-4h-4zM5 17h3.13c.02.09.06.17.09.25.24.68.65 1.28 1.18 1.75H5v-2zm14 2h-4.4c.54-.47.95-1.07 1.18-1.75.03-.08.07-.16.09-.25H19v2z")

		g = defsEl.append("symbol").attr("id", "ac_ic_node_container").attr("viewBox", "0 0 24 24").append("g")
			g.append("rect").attr("x","4").attr("y","2").attr("width","16").attr("height","20")
			g.append("path").attr("d", "M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5v-3h3.56c.69 1.19 1.97 2 3.45 2s2.75-.81 3.45-2H19v3zm0-5h-4.99c0 1.1-.9 2-2 2s-2-.9-2-2H5V5h14v9z")

		defsEl.append("symbol").attr("id", "ac_ic_node_cluster").attr("viewBox", "0 0 24 24").append("path").attr("d", "M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96z")

	}
	
	/**
	* 노드 랜더링
	*/
	private static renderNode(nodes:d3Select.Selection<SVGElement,any,SVGElement,any>)  {

		nodes
			.append("use")
			.attr("height", (d:model.Node) => {return (d.kind==model.NodeKind.NODE||d.kind==model.NodeKind.CLUSTER?40:30)} )
			.attr("width", (d:model.Node) => {return (d.kind==model.NodeKind.NODE||d.kind==model.NodeKind.CLUSTER?40:30)} )
			.attr("xlink:href", (d:model.Node) => `#ac_ic_node_${d.kind}`)
			.attr("x", (d:model.Node) => {return (d.kind==model.NodeKind.NODE||d.kind==model.NodeKind.CLUSTER?-20:-12)})
			.attr("y", (d:model.Node) => {return (d.kind==model.NodeKind.NODE||d.kind==model.NodeKind.CLUSTER?-20:-12)})

		// 라벨 Render
		nodes.append("text")
			.attr("class", "label")
			.text(d => { return d.name.length > 30 ? d.name.substring(0,27) + "...": d.name; })
			// .text(d => d.name)
			.attr("x", -24)
			.attr("y", 16)

	}
	/**
	* 노드 랜더링
	*/
	private static renderToolbar(svg:d3Select.Selection<SVGElement,any,SVGElement,any>)  {
	
	}

};	
