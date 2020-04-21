import * as d3					from "d3";
import * as model				from "./model";
import {KialiParser}			from "./kiali.parser";
import {ConfigModel}			from "../../model/models";
import {UI}						from "../../utils/lang";
import {GraphBase}				from "../graph.base";
import {EventBinder}			from "./event";
import {Traffic,GraphUI}		from "./lang";
import {LegendModel, Toolbar}	from "../toolbar";


require("./graph.css");

/**
 * 서비스 메시 그래프
 */
export class MeshGraph extends GraphBase {

	/**
	 * 그래프 그리는 단계
	 * 		1단계 그룹핑
	 * 			group-nodes 구성한다.
	 * 		2단계 : 그룹관 관게를 지정
	 * 			group-groups 구성한다.
	 * 		3단계 : 그룹내 nodes 들간 순서를 정의한다. 
	 * 			nodes depth 정의
	 * 		4단계 : 데이터를 화면에 그려준다.
	 * 
	 */
	public populate(conf:ConfigModel.Config, svgEl:d3.Selection<SVGSVGElement, any, SVGElement, any>, bounds:DOMRect, outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>) {

		let source:model.Source = KialiParser.parse(conf.data);
		if(!source) return;

		// 아이콘 defs:marker 그리기
		if(svgEl.select("defs").size() == 0) svgEl.append("defs").call(MeshGraph.renderDefs, conf);

		// outline 그리기
		outlineEl.call(MeshGraph.renderOutline, source, conf, bounds); 

		// SVG 클릭하면 selected 클리어 처리
		this.svg().call(EventBinder.onUnSelected, conf.mesh);

		// 범례
		const legends:Array<LegendModel> = [
			{
				header: "Node",
				rows: [
					{ label: "App", ico: '<g class="node app"><use class="ico" href="#ac_ic_node_app" width="20" height="20"></use></g>' },
					{ label: "Service", ico: '<g class="node app"><use class="ico" href="#ac_ic_node_service" width="20" height="20"></use></g>' },
					{ label: "Service Entry", ico: '<g class="node app"><use class="ico" href="#ac_ic_node_service_entry" width="20" height="20"></use></g>' }
				]
			},
			{
				header: "Node Colors",
				rows: [
					{ label: "Normal", ico: '<g class="node app"><use  class="ico app" href="#ac_ic_node_app" width="20" height="20"></use></g>' },
					{ label: "Warn", ico: '<g class="node app warn"><use class="ico app"  href="#ac_ic_node_app" width="20" height="20"></use></g>' },
					{ label: "Error", ico: '<g class="node app error"><use class="ico app"  href="#ac_ic_node_app" width="20" height="20"></use></g>' }
				]
			},
			{
				header: "Node Background",
				rows: [
					{ label: "Restrict Namespace", ico: '<g class="node app"><use  class="ico" href="#ac_ic_node_ns_restrict" width="20" height="20"></use></g>' },
					{ label: "External Namespace", ico: '<g class="node app"><use class="ico app"  href="#ac_ic_node_app_external" width="20" height="20"></use></g>' }
				]
			},
			{
				header: "Edges",
				rows: [
					{ label: "Warn", ico: '<g class="edge warn"><line x1="0" y1="10" x2="40" y2="10" marker-end="url(#ac_ic_arrowhead_warn)"></line></g>' },
					{ label: "Error", ico: '<g class="edge error"><line x1="0" y1="10" x2="40" y2="10" marker-end="url(#ac_ic_arrowhead_error)"></line></g>' },
					{ label: "Success", ico: '<g class="edge success"><line x1="0" y1="10" x2="40" y2="10" marker-end="url(#ac_ic_arrowhead_success)"></line></g>' },
					{ label: "Idle", ico: '<g class="edge"><line x1="0" y1="10" x2="40" y2="10" marker-end="url(#ac_ic_arrowhead_idle)"></line></g>' },
					{ label: "mTLS", ico: '<use href="#ac_ic_mtls" width="15" height="15"></use>' }
				]
			},
			{
				header: "Node Badges",
				rows: [
					{ label: "Missing Sidecar", ico: '<use href="#ac_ic_misssidecar" width="24" height="24"></use>' },
					{ label: "Virtual Services", ico: '<use href="#ac_ic_virtualservice" width="24" height="24"></use>' },
					{ label: "Circuit Breaker", ico: '<use href="#ac_ic_circuitbreaker" width="24" height="24"></use>' }
				]
			}
		]
		this.svg().call(Toolbar.render, this, legends);
	}

	/**
	 * outline 랜더링
	 * 
	 * @param outlineEl outline element
	 * @param source 데이터
	 * @param conf  옵션
	 * @param bounds  전체 크기
	 */
	private static renderOutline(outlineEl:d3.Selection<SVGGElement,model.Node,SVGElement,any>, source:model.Source, conf:ConfigModel.Config,  bounds:DOMRect) {

		if(!source || !source.groups) return;

		// #1 노드 랜더링
		outlineEl.call(MeshGraph.renderGroupsEl, 0, 0, source.groups, source.edges, conf, bounds)

		// #2 선 랜더링
		outlineEl.call(MeshGraph.renderLines, source.edges, conf, bounds);


	}


	/**
	 * g.groups  랜더링
	 * 
	 * @param parentEl outline 엘리먼트(ROOT) 이거나  부모 group 엘리먼트
	 * @param x X 좌표
	 * @param y Y 좌표
	 * @param groups 랜더링할 그룹 리스트
	 * @param mapEdges 노드 연결정보 
	 * @param conf  옵션
	 * @param bounds  전체 크기
	 */
	private static renderGroupsEl(parentEl:d3.Selection<SVGGElement,any,SVGElement,model.Group>, x:number, y:number, groups:Array<model.Group>, sourceEdge:model.EdgeSource, conf:ConfigModel.Config, bounds:DOMRect) {

		// group element 추가
		const groupsEl:d3.Selection<SVGGElement,any,SVGElement,any> = parentEl.append("g")
			.attr("class","groups")

		groupsEl.selectAll("g.group")
			.data(groups).enter()
			.append("g")
			.attr("class", "group")
			.each((d:model.Group, i:number, els:SVGGElement[])=> {
				d3.select(els[i]).call(MeshGraph.renderGroupEl, x, y, sourceEdge, conf, bounds);
			})

		// 부모children 그리고 난 후  보모 엘리먼트내에서 수평정렬
		UI.alignVertical(groupsEl.node());

	}

	/**
	 * g.group 내  g.nodes 노드들 랜더링하고 하위 children 이 있다면 g.groups 노드 랜더링
	 * 
	 * @param groupEl group 엘리먼트
	 * @param x X 좌표
	 * @param y Y 좌표
	 * @param mapEdges 노드 연결정보 
	 * @param conf  옵션
	 * @param bounds  전체 크기
	 */
	private static renderGroupEl(groupEl:d3.Selection<SVGGElement,any,SVGElement,any>, x:number,y:number, sourceEdge:model.EdgeSource, conf:ConfigModel.Config, bounds:DOMRect) {
		
		let group:model.Group = groupEl.datum();
		let isOutline:boolean = (conf.mesh.type=="service"||conf.mesh.type=="workload") ? false: (group.nodes.length > 1);	// isOutline=false & 노드가 1개인 경우(dGroup.nodes.length==1) => outline 없도록 처리

		// 위치 조정
		groupEl.attr("transform", (d:model.Group, i:number, els:SVGGElement[]) => {
			let h:number = els[i].parentElement.getBoundingClientRect().height;
			y = h==0?0: h+conf.mesh.group.distance.y;	//height 0이면 그룹에서 children 에서 첫번째 그룹임
			return `translate(${x} ${y})`;
		})

		// group element 추가
		let nodesEl:d3.Selection<SVGGElement, any, SVGElement, any> = groupEl.append("g")
			.attr("class", (d:model.Group) => { 
				// 에러비율에 따른 css 선택 - service는 제외하고 포함하고 있는 노드들의 max 에러비율
				let err = 0;
				d.nodes.forEach( (el:model.Node)=> {
					//if(el.nodeType != "service" && el.health) err = Math.max(err, el.health.requests.errorRatio);
					if(el.health) err = Math.max(err, el.health.requests.errorRatio);
				});
				if(err > conf.global.health.error.ratio) return "nodes error";
				else if(err > conf.global.health.warn.ratio) return "nodes warn" ;
				else return "nodes";
			})
			.call(EventBinder.onHighlight, sourceEdge, conf.mesh)
			.call(EventBinder.onSelected, sourceEdge, conf.mesh)

		// 노드  데이터가 있다면 하위에 nodes 추가
		if (group.nodes.length > 0) {

			// g.nodes 추가
			nodesEl.call(MeshGraph.renderNodesEl, sourceEdge, conf)

			if(isOutline) {

				//  group-boundary BOX 추가
				const rect:DOMRect = (<SVGGElement>nodesEl.node()).getBBox();

				// BOX 크기는 padding config 적용
				const h:number = rect.height + conf.mesh.group.padding.top + conf.mesh.group.padding.bottom;	//높이
				const w:number = rect.width + conf.mesh.group.padding.left + conf.mesh.group.padding.right;		//너비
				const x:number = rect.x - conf.mesh.group.padding.left;
				const y:number = rect.y - conf.mesh.group.padding.top;

				nodesEl.insert("rect", "g.node:first-child")
					.attr("x", x)
					.attr("y", y)
					.attr("width", w)
					.attr("height", h)

				// 제목 추가
				if (conf.mesh.group.label.show) {
					nodesEl.append("text")
						.attr("class", "caption")
						.text(group.name)
						.attr("x", (d:any, idx:number, els:SVGTextElement[]) => {
							return x + (w - els[idx].getBBox().width)/2;	//글짜 width 에 따른 가운데 정렬
						})
						.attr("y", h)
				}

			}

			// group 내에 다시 children (groups) 가 있다면
			if(group.groups && group.groups.length > 0) {

				// children(groups) 추가
				groupEl.call(MeshGraph.renderGroupsEl, nodesEl.node().getBoundingClientRect().width + conf.mesh.node.distance.x, 0, group.groups, sourceEdge, conf, bounds);

				// children 이 2개 이상이면 현재 그룹내 g.nodes 를 수직 정렬
				if(group.groups.length > 1) UI.alignVertical(<SVGElement>groupEl.select("g.nodes").node());
			}

		}

	}	


	/**
	 * g.nodes 엘리먼트 하위 g.node 엘리먼트들(service, app)을 랜더링
	 * 
	 * @param nodesEl g.nodes 엘리먼트 
	 * @param mapEdges 노드 연결정보 
	 * @param conf  옵션
	 */
	private static renderNodesEl(nodesEl:d3.Selection<SVGGElement,any,SVGElement,any>, sourceEdge:model.EdgeSource, conf:ConfigModel.Config) {

		nodesEl.selectAll("g.node")
			.data(nodesEl.datum().nodes).enter()
			.append("g")
			.attr("class", (d:model.Node) => { 
				// 에러비율에 따른 css 선택
				if(!d.health || !conf.mesh.alert.show) return `node ${d.nodeType}`;
				else if(d.health.requests.errorRatio >= conf.global.health.error.ratio) return conf.mesh.alert.error.animation.node ? `node ${d.nodeType} error anmi`: `node ${d.nodeType} error`;
				else if(d.health.requests.errorRatio >= conf.global.health.warn.ratio) return `node ${d.nodeType} warn`;
				else return `node ${d.nodeType}`;
			})
			.attr("id", (d:model.Node) => d.id)
			.call(EventBinder.onHighlight, sourceEdge, conf.mesh)
			.call(EventBinder.onSelected, sourceEdge, conf.mesh)
			.each((d:model.Group, i:number, els:SVGGElement[])=> {
				d3.select(els[i]).call(MeshGraph.renderNodeEl, conf);
			})
			.call(MeshGraph.alignNodes, conf)

	}

	/**
	 * 1개의 g.node 엘리먼트의 라벨, 뱃지, 아이콘 추가
	 * 
	 * @param nodeEl g.node 엘리먼트
	 * @param conf  옵션
	 */
	private static renderNodeEl(nodeEl:d3.Selection<SVGGElement,model.Node,null,any>, conf:ConfigModel.Config) {

		const isLabel:boolean = conf.mesh.node.label.show;
		const radius:number = conf.mesh.node.r;

		// 라벨그룹 추가
		let labelEl:d3.Selection<SVGGElement,any,SVGElement,any> = nodeEl.append("g")
			.attr("class","label")
		
		// node element 별 라벨 추가
		let txtEl:d3.Selection<SVGGElement,any,SVGElement,any> = labelEl.append("text")
			.text((d:model.Node) => {
				if(!isLabel) return  " ";
				else if(d.nodeType=="service") return d.service;
				else if(d.version && d.version !="unknown") return d.version;
				else return d.workload?d.workload:d.app;
			})
		
		// 문자 높이 (기준값)
		const h:number = txtEl.node().getBoundingClientRect().height;

		// 라벨그룹 뱃지 추가
		if(nodeEl.datum().hasVS || nodeEl.datum().hasMissingSC || nodeEl.datum().hasCB) {
			const iconHW:number = h * 1.2	// 아이콘 너비,높이 조정값
			let shiftX:number = 0;			// text 좌측이동 좌표

			if(nodeEl.datum().hasMissingSC)  {
				labelEl.append("use")
					.attr("height", iconHW).attr("width", iconHW)
					.attr("x", shiftX)
					.attr("xlink:href", "#ac_ic_misssidecar");
				shiftX += h * 0.8;
			}
			if(nodeEl.datum().hasVS)  {
				labelEl.append("use")
					.attr("height", iconHW).attr("width", iconHW)
					.attr("x", shiftX)
					.attr("xlink:href", "#ac_ic_virtualservice");
				shiftX += h * 0.8;
			}

			if(nodeEl.datum().hasCB)  {
				labelEl.append("use")
					.attr("height", iconHW).attr("width", iconHW)
					.attr("x", shiftX)
					.attr("xlink:href", "#ac_ic_circuitbreaker");
				shiftX += h * 0.8;
			}

			// 텍스트 뱃지크기 만큼 우측으로 이동 (0.8 배 shift 되었음.. 라벨과는 1배 유격을 위해 0.2배 추가)
			txtEl.attr("x", shiftX + (h * 0.2));
		}


		// 아이콘 추가
		const w:number = labelEl.node().getBoundingClientRect().width
		nodeEl.append("use")
			.attr("xlink:href", (d:model.Node) => {
				if(d.nodeType=="service") {
					if(d.isInaccessible) return "#ac_ic_node_ns_restrict";
					else if (d.isServiceEntry == "MESH_EXTERNAL") return "#ac_ic_node_service_entry";
					else  return d.isOutside? "#ac_ic_node_service_external": "#ac_ic_node_service";
				} else {
					return d.isOutside ? "#ac_ic_node_app_external": "#ac_ic_node_app";
				}
			})
			.attr("x", w/2 - radius).attr("y", h)
			.attr("class", d=> `ico ${d.nodeType}`)

	}	

	/**
	 * g.nodes 엘리먼트 내  g.node 엘리먼트들 간 위치 조정
	 * 
	 * @param nodeEl g.nodes 엘리먼트
	 * @param conf  옵션
	 */
	private static alignNodes(nodeEl:d3.Selection<SVGGElement,any,SVGElement,any>, conf:ConfigModel.Config) {

		const serviceEl:d3.Selection<SVGGElement,any,SVGElement,any> = nodeEl.filter((d:model.Node)=>d.nodeType=="service")
		const appEl:d3.Selection<SVGGElement,any,SVGElement,any> = nodeEl.filter((d:model.Node)=>d.nodeType=="app"||d.nodeType=="workload")

		// elemets 크기 변경으로 인한 위치 조정
		let cursor = { 
			offsetX:0, 		// 시작 x offset (x 조정값)
			x:0, 			// 현재 x
			y:0, 			// 현재 y
			maxWidth: 0		// 최대너비
		};
		
		// 첫번째 섹션 (service)
		serviceEl
			.each( (d:model.Node, idx:number, els:SVGGElement[]) => {
				let box:DOMRect = els[idx].getBBox();
				// x 시작지점과 maxWidth 값을 얻는다 (maxwidth 값이 맞춰주려고함)
				cursor.offsetX = Math.max(-box.x, cursor.offsetX);
				cursor.maxWidth = Math.max(box.width, cursor.maxWidth);
			})
			.attr("transform", (d:model.Node, idx:number, els:SVGGElement[]) => { 
				let box = els[idx].getBBox(), y = cursor.y;
				cursor.y += (box.height + conf.mesh.node.distance.y); 			// y 축 좌표 계산
				return `translate(${cursor.x+cursor.offsetX} ${y})`;
			})

		// 서비스 존재하면 - 첫번재와 두번째 섹션은 우측으로 이동(distance 적용)
		if(serviceEl.size() > 0) {
			cursor.x += cursor.maxWidth + conf.mesh.node.distance.x;
			cursor.y = 0;
			cursor.maxWidth = 0;
		}
		
		// 두번째 섹션 (app)
		appEl
			.each( (d:model.Node, idx:number, els:SVGGElement[]) => {
				let box:DOMRect = els[idx].getBBox();
				cursor.offsetX = Math.max(box.x*-1, cursor.offsetX);
				cursor.maxWidth = Math.max(box.width, cursor.maxWidth);
			})
			.attr("transform", (d:model.Node, idx:number, els:SVGGElement[]) => { 
				let box:DOMRect = els[idx].getBBox(), y = cursor.y;
				cursor.y += box.height + conf.mesh.node.distance.y;
				return `translate(${cursor.x+cursor.offsetX} ${y})`;
			})
		
		if(serviceEl.size() == 1 && appEl.size() > 1)  UI.alignVertical(serviceEl.node());	// 서비스 1개 + app 이 n개가 있으는 경우 service 수직정렬
		else if(serviceEl.size() > 1 && appEl.size() == 1)  UI.alignVertical(appEl.node());	
			
	}

	/**
	 * 라인 랜더링
	 * 
	 * @param outlineEl outline 엘리먼트
	 * @param edges 라인 정보(전체)
	 * @param mapEdges 노드 연결정보 
	 * @param conf  옵션
	 */
	private static renderLines(outlineEl:d3.Selection<SVGGElement,any,SVGElement,any>, sourceEdge:model.EdgeSource, conf:ConfigModel.Config, bounds:DOMRect) {


		GraphUI.calcEdgesPosition(outlineEl, sourceEdge, conf.mesh, bounds); 

		// 선 그룹
		let g:d3.Selection<SVGGElement,any,SVGElement,any> = outlineEl.append("g")
			.attr("class","edges")
			.selectAll("g.edge")
				.data(sourceEdge.data)
				.enter()
				.append("g")
				.attr("id",d=>d.id)
				.attr("class", (d:model.Edge) => conf.mesh.alert.show ? `edge ${Traffic.getStatus(d, conf)}`: "edge")	//에러비율에 따른 라인 css 선택
				.call(EventBinder.onHighlight, sourceEdge, conf.mesh)
				.call(EventBinder.onSelected, sourceEdge, conf.mesh)

		// 선 
		g.insert("line")
			.attr("x1",d=>d.x1).attr("y1",d=>d.y1)
			.attr("x2",d=>d.x2).attr("y2",d=>d.y2)
			.attr("marker-start", (d:model.Edge) => d.isTwoWay?`url(#ac_ic_arrowhead_${Traffic.getStatus(d, conf)}${d.isTwoWay?"_r":""} )`:"")
			.attr("marker-end", (d:model.Edge) => `url(#ac_ic_arrowhead_${Traffic.getStatus(d, conf)})`);
		
		// 라벨
		if(conf.mesh.edge.label.type != "none") {
			g.each((d:model.Edge, i:number, els:SVGGElement[])=> {
					d3.select(els[i]).call(MeshGraph.renderLineLabel, conf);
				})
		}


		// 트래픽 에러에 대한 Alert(error ratio)
		if(conf.mesh.alert.show && conf.mesh.alert.error.label.show) {
			g.filter(".error")
				.append("text")
				.attr("class", "label error")
				.text((d:model.Edge)=> {
					return `${100-Traffic.getErrorRatio(d.traffic)} %`;
				})
				.attr("x", (d:any, i:number, els:SVGTextElement[])=> {
					return d.x1+(d.x2-d.x1)/2 - (els[i].getBBox().width/2);
				})
				.attr("y", (d:any, i:number, els:SVGTextElement[])=> {
					return d.y1 + (d.y2-d.y1)/2 + 5;	//선과 라벨과의 거리
				})
		}

		// 트래픽 애니메이션
		if(conf.mesh.edge.traffic.animation.show) {
			let duration:number = conf.mesh.edge.traffic.animation.duration;

			// 단방향 #1
			g.filter(d =>!d.isTwoWay &&d.traffic.responses)
				.append("circle")
				.attr("r", 4)
				.attr("class","anmi")
				.attr("transform", (d:any) => {
					return `translate(${d.x1},${d.y1})`;
				})
				.transition()
				.duration(0)
				.on("start", function repeat() {
					d3.active(this)
						.transition().duration(duration/2 *Math.floor(Math.random() * 3)).attr('transform', (d:any)=>`translate(${d.x2},${d.y2})`)
						.transition().duration(0).attr("transform", (d:any)=>`translate(${d.x1},${d.y1})`)
						.on("start", repeat);
				});

			// 양방향 #1
			g.filter(d =>d.isTwoWay && d.traffic.responses)
				.append("circle")
				.attr("r", 4)
				.attr("class","anmi")
				.attr("transform", (d:any) => {
					return `translate(${d.x1},${d.y1-4})`;
				})
				.transition().duration(0)
				.on("start", function repeat() {
					d3.active(this)
						.transition().duration(duration/2 *Math.floor(Math.random() * 3)).attr('transform', (d:any)=>`translate(${d.x2},${d.y2-4})`)
						.transition().duration(0).attr("transform", (d:any)=>`translate(${d.x1},${d.y1-4})`)
						.on("start", repeat);
				});

			// 양방향 #2
			g.filter(d =>d.isTwoWay && d.traffic.responses)
				.append("circle")
				.attr("r", 4)
				.attr("class","anmi")
				.attr("transform", (d:any) => {
					return `translate(${d.x2},${d.y2+4})`;
				})
				.transition().duration(0)
				.on("start", function repeat() {
					d3.active(this)
						.transition().duration(duration/2 *Math.floor(Math.random() * 3)).attr('transform', (d:any)=>`translate(${d.x1},${d.y1+4})`)
						.transition().duration(0).attr("transform", (d:any)=>`translate(${d.x2},${d.y2+4})`)
						.on("start", repeat);
				});

		}

	}

	/**
	 * 라인 라벨 랜더링
	 * 
	 * @param edgeEl g.edge 엘리먼트
	 * @param conf  옵션
	 */
	private static renderLineLabel(edgeEl:d3.Selection<SVGGElement,any,SVGElement,any>, conf:ConfigModel.Config) {

		const edge:model.Edge = edgeEl.datum();

		if(conf.mesh.edge.label.type == "responseTime" || conf.mesh.edge.label.type == "requestsPercent" || edge.isMTLS) {

			let g:d3.Selection<SVGGElement,any,SVGElement,any> = edgeEl.append("g")
				.attr("class","label");

			// 문자열
			if(conf.mesh.edge.label.type == "responseTime" && edge.responseTime) {
				g.append("text").text( `${edge.responseTime}ms`);
			} else if(conf.mesh.edge.label.type == "requestsPercent" && edge.traffic && edge.traffic.rates && edge.traffic.rates.httpPercentReq) {
				g.append("text").text(d=> `${d.traffic.rates.httpPercentReq}%`);
			} else  {
				g.append("text").text(d=> `_`).attr("visibility","hidden");	// 아이콘 사이즈 계산을 위해 임의 문자 삽입후 hidden
			}

			const iconWH = g.node().getBoundingClientRect().height;

			// mTSL 뱃지
			let shiftX:number = 0;
			if(edge.isMTLS) {
				g.append("use")
					.attr("height", iconWH)
					.attr("width", iconWH)
					.attr("xlink:href", "#ac_ic_mtls")
					.attr("transform", `translate(${-iconWH},-2)`);
				shiftX  = iconWH;
			}

			//g.label 위치 이동
			//		x:	선 중간에서 너비(1/2)를 반영
			//			badge 가 있는 경우는 좌측으로 너비만큼 이동 했으므로 x - shift  반영)
			//		y:	선  중간에서 높이 적용(-)
			const x = edge.x2 - (edge.x2-edge.x1)/2	- g.node().getBoundingClientRect().width/2 + shiftX;	
			const y = edge.y2 - (edge.y2-edge.y1)/2 - iconWH;
			g.attr("transform", `translate(${x} ${y})`)

		}

	}

	/**
	 * defs 정의
	 * 
	 * @param defsEl def 엘리먼트
	 */
	private static renderDefs(defsEl:d3.Selection<SVGDefsElement, any, SVGElement, any>, conf:ConfigModel.Config) {

		const radius:number = conf.mesh.node.r;

		// ac_ic_arrowhead - arrow head
		defsEl.html(
`<marker id="ac_ic_arrowhead_idle" refX="6" refY="6" markerWidth="20" markerHeight="20" markerUnits="userSpaceOnUse" orient="auto"><path d="M 0 0 L 8 6 L 0 12  L 2 6"></path></marker>
<marker id="ac_ic_arrowhead_warn" refX="6" refY="6" markerWidth="20" markerHeight="20" markerUnits="userSpaceOnUse" orient="auto"><path d="M 0 0 L 8 6 L 0 12  L 2 6"></path></marker>
<marker id="ac_ic_arrowhead_error" refX="6" refY="6" markerWidth="20" markerHeight="20" markerUnits="userSpaceOnUse" orient="auto"><path d="M 0 0 L 8 6 L 0 12  L 2 6"></path></marker>
<marker id="ac_ic_arrowhead_success" refX="6" refY="6" markerWidth="20" markerHeight="20" markerUnits="userSpaceOnUse" orient="auto"><path d="M 0 0 L 8 6 L 0 12  L 2 6"></path></marker>
<marker id="ac_ic_arrowhead_success_r" refX="6" refY="6" markerWidth="20" markerHeight="20" markerUnits="userSpaceOnUse" orient="auto"><path d="M 4 8 L 16 0 L 12 6 L 16 12 L 4 8"></path></marker>
`);

		// ac_ic_node_service - 노드 ICON - service
		// ac_ic_node_app - 노드 ICON - app
		// ac_ic_node_ns_restrict / 노드 ICON
		// ac_ic_node_app_external / 노드 ICON
		// ac_ic_node_service_external / 노드 ICON 
		// ac_ic_node_service_entry / 노드 ICON - Service Entry 
		// ac_ic_mtls / 선 라벨 ICON - mTLS
		// ac_ic_misssidecar / missing sidecar
		// ac_ic_virtualservice / virtual service
		// ac_ic_circuitbreaker / circuitbreaker

		defsEl.append("symbol").html(`<svg id="ac_ic_node_service" viewBox="0 0 30 30" width="${radius*2}" height="${radius*2}"><polygon points="15,0  30,30  0, 30"></polygon></svg>`)
		defsEl.append("symbol").html(`<svg id="ac_ic_node_app" viewBox="0 0 30 30" width="${radius*2}" height="${radius*2}"><circle cy="15" cx="15" cy="14" r="14"></circle></svg>`)
		defsEl.append("symbol").html(`<svg id="ac_ic_node_service_entry" viewBox="0 0 30 30" width="${radius*2}" height="${radius*2}"><polygon points="1,1  22,1  30,15 22,30 1,30"></polygon></svg>`)
		defsEl.append("symbol").html(`<svg id="ac_ic_node_app_external" viewBox="0 0 30 30" width="${radius*2}" height="${radius*2}"><path d="M15,1 a14,14 0 1,0 0.00001,0"></path><path d="M19.21 6.76l4 4v-4zm4 12l-4 4h4zm-12 4l-4-4v4zm-4-12l4-4h-4zm12.95-.95c-2.73-2.73-7.17-2.73-9.9 0s-2.73 7.17 0 9.9 7.17 2.73 9.9 0 2.73-7.16 0-9.9zm-1.1 8.8c-2.13 2.13-5.57 2.13-7.7 0s-2.13-5.57 0-7.7 5.57-2.13 7.7 0 2.13 5.57 0 7.7z"></path></svg>`);
		defsEl.append("symbol").html(`<svg id="ac_ic_node_service_external" viewBox="0 0 30 30" width="${radius*2}" height="${radius*2}"><polygon points="${radius},0  ${radius*2},${radius*2}  0, ${radius*2}"></polygon><path d="M19.25 10.76l4 4v-4zm4 12l-4 4h4zm-12 4l-4-4v4zm-4-12l4-4h-4zm12.95-.95c-2.73-2.73-7.17-2.73-9.9 0s-2.73 7.17 0 9.9 7.17 2.73 9.9 0 2.73-7.16 0-9.9zm-1.1 8.8c-2.13 2.13-5.57 2.13-7.7 0s-2.13-5.57 0-7.7 5.57-2.13 7.7 0 2.13 5.57 0 7.7z"></path></svg>`);
		defsEl.append("symbol").html(`<svg id="ac_ic_node_ns_restrict" viewBox="0 0 30 30" width="${radius*2}" height="${radius*2}"><path d="M15,0 L29,29 L0,29 L15,0"></path><path transform="translate(2,7)" d="M12.65 10C11.83 7.67 9.61 6 7 6c-3.31 0-6 2.69-6 6s2.69 6 6 6c2.61 0 4.83-1.67 5.65-4H17v4h4v-4h2v-4H12.65zM7 14c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2z"></path></svg>`)
		defsEl.append("symbol").attr("id", "ac_ic_mtls").attr("viewBox", "0 0 24 24").append("path").attr("d", "M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z")
		defsEl.append("symbol").attr("id", "ac_ic_misssidecar").attr("viewBox", "0 0 30 30").append("path").attr("d", "M3 13h2v-2H3v2zm0 4h2v-2H3v2zm2 4v-2H3c0 1.1.89 2 2 2zM3 9h2V7H3v2zm12 12h2v-2h-2v2zm4-18H9c-1.11 0-2 .9-2 2v10c0 1.1.89 2 2 2h10c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 12H9V5h10v10zm-8 6h2v-2h-2v2zm-4 0h2v-2H7v2z")
		defsEl.append("symbol").attr("id", "ac_ic_virtualservice").attr("viewBox", "0 0 30 30").append("path").attr("d","M17 16l-4-4V8.82C14.16 8.4 15 7.3 15 6c0-1.66-1.34-3-3-3S9 4.34 9 6c0 1.3.84 2.4 2 2.82V12l-4 4H3v5h5v-3.05l4-4.2 4 4.2V21h5v-5h-4z")
		defsEl.append("symbol").attr("id", "ac_ic_circuitbreaker").attr("viewBox", "0 0 30 30").append("path").attr("d","M7 2v11h3v9l7-12h-4l4-8z")
		defsEl.append("symbol").attr("id", "ac_ic_close").attr("viewBox", "0 0 24 24").append("path").attr("d","M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z")


	}

};	
