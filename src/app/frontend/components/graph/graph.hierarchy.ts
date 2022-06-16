"uss strict"
import * as d3						from "d3";
import * as d3Select				from "d3-selection";
import {HierarchyModel as model}	from "@/components/graph/model/graph.model";
import {Config}						from "@/components/graph/model/config.model";
import {UI, WH}						from "@/components/graph/utils/ui";
import {GraphBase}					from "@/components/graph/graph.base";
import "@/components/graph/graph.hierarchy.css";

/**
 * Topology 그래프 랜더러
 */
export default class HierarchyGraph extends GraphBase {


	/**
	 * (abstract) 랜더링
	 * 
	 * @param outlineEl 외곽 g element (zoom & drag 용)
	 * @param bounds 랜더링 영역 크기 (x,y,height,width)
	 * @param conf 데이터 & 옵션
	 */
	public populate(outlineEl:d3Select.Selection<SVGGElement,any,SVGElement,any>, bounds:WH, conf:Config) {
		
		// Set width
		let width:number = bounds.width;
		if(conf.extends.hierarchy.scale.minWidth > 0 && bounds.width < conf.extends.hierarchy.scale.minWidth*conf.global.scale.ratio)  width = conf.extends.hierarchy.scale.minWidth;
		if(conf.extends.hierarchy.scale.maxWidth > 0 && bounds.width > conf.extends.hierarchy.scale.maxWidth*conf.global.scale.ratio)  width = conf.extends.hierarchy.scale.maxWidth;
		width -= (conf.extends.hierarchy.group.box.border.width*2);	 // border

		// svg > defs
		if(this.svg.select("defs").size() == 0) this.svg.append("defs").call(HierarchyGraph.renderDefs, conf);

		// data 가공
		let data:Array<model.Node> = [];
		Object.keys(conf.data).forEach( (k:string)=> {
			let d:Array<model.Node> = conf.data[k];
			const root = d.reduce((acc, cur:model.Node) => {
				if(cur.ownerReference && cur.ownerReference.kind && cur.ownerReference.name) {
					d.reduce((a:model.Node, c:model.Node) => {
						if(c.kind == cur.ownerReference!.kind && c.name == cur.ownerReference!.name) {
							if(!c.children) c.children=[]
							c.children.push(cur)
						}
						return a
					}, new model.Node());
				} else {
					if(!cur.children) cur.children = [];
					acc.children.push(cur)
				}
				return acc;
			}, new model.Node(k))
			data.push(root)
		});
		// rendering groups
		// svg > g.graph > g.outlineWrap > g.outline > g.group
		//		> text
		//      > g.boxWrap > g.box > g.tree
		let gY = 0;
		const padding = conf.extends.hierarchy.group.box.padding;
		const treeWidth:number = width - (padding.left + padding.right);

		data.forEach((d:model.Node)=> {
			
			const g:d3.Selection<SVGGElement, any, SVGElement, any> = outlineEl.append("g").attr("class","group");
			let t; 
			if (conf.extends.hierarchy.group.title.display == "always" || conf.extends.hierarchy.group.title.display == "has" &&  d.children.length > 0) {
				t = g.append("text").text(d.name).attr("transform", (d:any,i:number,els:SVGTextElement[]|d3.ArrayLike<SVGTextElement>)=> {
					return `translate(0,${els[i].getBBox().y * -1})`
				})
			}

			if(d.children.length > 0) {
				let h = t ? t.node()!.getBBox().height + conf.extends.hierarchy.group.title.spacing:0;
				UI.appendBox(g, (box: d3.Selection<SVGGElement, any, SVGElement, any>)=> {
					d.children.forEach((c:model.Node)=> {
						let gg = box.append("g").attr("class","tree")
							.call(HierarchyGraph.renderHierarchy, c, conf, treeWidth)
							.attr("transform", (d:any,i:number,els: Array<SVGGElement>|d3.ArrayLike<SVGGElement>)=> {
								return `translate(0,${h-els[i].getBBox().y})`
							});
						h += gg.node()!.getBBox().height + conf.extends.hierarchy.group.box.tree.spacing;	// multi-root 간 간격
					});
				}, width, padding, conf.extends.hierarchy.group.box.background, conf.extends.hierarchy.group.box.border);

			}
			// + move XY
			g.attr("transform", `translate(${(bounds.width-width)/2},${gY})`)
			if(d.children.length > 0) gY += g.node()!.getBBox().height + conf.extends.hierarchy.group.spacing;
		});

		// toolbar aline default 값 정의 -  "none"(사용자 지정 X)이면
		if(conf.global.toolbar.align.horizontal == "none") conf.global.toolbar.align.horizontal = "right";
		if(conf.global.toolbar.align.vertical == "none") conf.global.toolbar.align.vertical = "top";


	}

	/**
	 * Hierarchy(tree) 랜더링 
	 * 
	 * @param data  랜더링 데이터
	 * @param treeWidth 너비 - 각 노드 너비 계산
	*/
	private static renderHierarchy(parentEl:d3Select.Selection<SVGGElement,any,SVGElement,any>, data:model.Node, conf:Config, treeWidth:number) {

		const nodeHeight:number = conf.extends.hierarchy.group.box.tree.node.height;	//default:30
		const nodeWidth:number = treeWidth/ 3;
		const icoWH:number = nodeHeight-2;
		const marginW:number = 2.5;	// margin(2.5) - between icon and text, between text and text

		const layoaut = d3.tree().nodeSize([nodeHeight, nodeWidth]);

		let d:d3.HierarchyNode<model.Node> = d3.hierarchy(data, (d:any) => d.children);	//  assigns the data to a hierarchy using parent-child relationships
		let nodes:d3.HierarchyPointNode<model.Node> = <d3.HierarchyPointNode<model.Node>>layoaut(d) // maps the node data to the tree layout

		// x-> y, y->x (because horizontal)
		nodes.each( (nd:d3.HierarchyPointNode<model.Node>)=> {
			if(nd.data.depth > 0) {
				nd.y =  nodeWidth * nd.data.depth;
			} else {
				nd.y =  nodeWidth * nd.depth;
			}
		})

		// adds each node as a group
		let nodeEl:d3.Selection<SVGGElement, any, SVGElement, any> = parentEl.selectAll(".node")
			.data(nodes.descendants())
		.enter().append("g")
			.attr("class", (conf.on && conf.on.nodeclick)? "node click": "node")
			.attr("transform", (d:d3.HierarchyPointNode<model.Node>) => `translate(${d.y},${d.x})`);

		// on nodeclick
		if(conf.on && conf.on.nodeclick) nodeEl.on("click", conf.on.nodeclick);

		// adds the icon to the node
		nodeEl.append("use")
			.attr("class","ico").attr("height",icoWH).attr("width",icoWH)
			.attr("xlink:href", (d:d3.HierarchyPointNode<model.Node>)=>`#ac_ic_${(d.data.kind || "").toLowerCase()}`)

		nodeEl.append("text")
			.text((d:d3.HierarchyPointNode<model.Node>) =>d.data.name)
			.attr("x", icoWH + marginW)
			.attr("y", (d:any,i:number,els: Array<SVGTextElement>|d3.ArrayLike<SVGTextElement>) => {
				const box = els[i].getBBox();
				return (nodeHeight-box.height)/2 - box.y;	//set vertical-middle
			})
			.each( (d:any,i:number,els:SVGTextElement[]|d3.ArrayLike<SVGTextElement>) =>{
				d.width = UI.ellipsisText(els[i], nodeWidth);	//calculate - text width
			});

		// adds the links between the nodes
		parentEl.selectAll(".link")
			.data( nodes.descendants().slice(1))
		.enter().append("path")
			.attr("class", "link")
			.attr("d", (d:any) => {
				// x->y, y->x (because horizontal)
				const x1 = d.parent!.y + d.parent!.width + icoWH + (marginW*2);
				const y1 = d.parent!.x + nodeHeight/2; // node height/2;
				const x2 = d.y;
				const y2 = d.x + nodeHeight/2; // node height/2;
				return `M ${x2},${y2} C ${(x2+x1)/2},${y2} ${(x2+x1)/2},${y1} ${x1},${y1}`;
			})

	}

	/**
	 * defs 정의
	 * 
	 * @param defsEl def 엘리먼트
	 */
	private static renderDefs(defsEl:d3.Selection<SVGDefsElement, any, SVGElement, any>) {

		// https://github.com/kubernetes/community/tree/master/icons
		defsEl.append("symbol").attr("id", "ac_ic_namespace")
			.attr("width", "18.035334mm").attr("height", "17.500378mm").attr("viewBox", "0 0 18.035334 17.500378")
			.html(`<g transform="translate(-0.99262638,-1.174181)">
<g transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">
	<path d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z" style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="M -6.8523435,3.8176372 A 1.1814304,1.171762 0 0 0 -7.3044284,3.932904 l -6.1787426,2.9512758 a 1.1814304,1.171762 0 0 0 -0.639206,0.794891 l -1.523915,6.6308282 a 1.1814304,1.171762 0 0 0 0.160175,0.89893 1.1814304,1.171762 0 0 0 0.06736,0.09281 l 4.276094,5.317236 a 1.1814304,1.171762 0 0 0 0.92363,0.440858 l 6.8576188,-0.0015 a 1.1814304,1.171762 0 0 0 0.9236308,-0.44011 l 4.2745966,-5.317985 a 1.1814304,1.171762 0 0 0 0.228288,-0.990993 L 0.53894439,7.6775738 A 1.1814304,1.171762 0 0 0 -0.10026101,6.8834313 L -6.2790037,3.9321555 A 1.1814304,1.171762 0 0 0 -6.8523435,3.8176372 Z m 0.00299,0.4550789 a 1.1191255,1.1099671 0 0 1 0.5426517,0.1085303 l 5.85315169,2.7948425 A 1.1191255,1.1099671 0 0 1 0.15197811,7.9290648 L 1.598051,14.21035 a 1.1191255,1.1099671 0 0 1 -0.2163123,0.939348 l -4.0493032,5.037304 a 1.1191255,1.1099671 0 0 1 -0.8749789,0.416906 l -6.4961006,0.0015 a 1.1191255,1.1099671 0 0 1 -0.874979,-0.417652 l -4.0508,-5.036554 a 1.1191255,1.1099671 0 0 1 -0.06362,-0.08832 1.1191255,1.1099671 0 0 1 -0.151942,-0.851028 l 1.443827,-6.2812853 a 1.1191255,1.1099671 0 0 1 0.605524,-0.7529758 l 5.8524036,-2.7963395 a 1.1191255,1.1099671 0 0 1 0.4288819,-0.1085303 z" style="color:#000000;font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:medium;line-height:normal;font-family:Sans;-inkscape-font-specification:Sans;text-indent:0;text-align:start;text-decoration:none;text-decoration-line:none;letter-spacing:normal;word-spacing:normal;text-transform:none;writing-mode:lr-tb;direction:ltr;baseline-shift:baseline;text-anchor:start;display:inline;overflow:visible;visibility:visible;fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;marker:none;enable-background:accumulate"/>
</g>
<text y="16.811775" x="9.9717083" style="font-style:normal;font-weight:normal;font-size:10.58333302px;line-height:6.61458349px;font-family:Sans;letter-spacing:0px;word-spacing:0px;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1">
	<tspan y="16.811775" x="9.9717083" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:2.82222223px;font-family:Arial;-inkscape-font-specification:'Arial, Normal';text-align:center;writing-mode:lr-tb;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.26458332px">ns</tspan>
</text>
<rect y="6.3689628" x="6.1734986" height="6.6900792" width="7.6735892" style="opacity:1;fill:none;fill-opacity:1;fill-rule:nonzero;stroke:#ffffff;stroke-width:0.40000001;stroke-linecap:butt;stroke-linejoin:round;stroke-miterlimit:10;stroke-dasharray:0.80000001, 0.4;stroke-dashoffset:3.44000006;stroke-opacity:1" />
</g>`)


		defsEl.append("symbol").attr("id", "ac_ic_deployment")
			.attr("width", "18.035334mm").attr("height", "17.500378mm").attr("viewBox", "0 0 18.035334 17.500378")
			.html(`<g transform="translate(-0.99262638,-1.174181)">
<g transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">
	<path d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z" style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="M -6.8523435,3.8176372 A 1.1814304,1.171762 0 0 0 -7.3044284,3.932904 l -6.1787426,2.9512758 a 1.1814304,1.171762 0 0 0 -0.639206,0.794891 l -1.523915,6.6308282 a 1.1814304,1.171762 0 0 0 0.160175,0.89893 1.1814304,1.171762 0 0 0 0.06736,0.09281 l 4.276094,5.317236 a 1.1814304,1.171762 0 0 0 0.92363,0.440858 l 6.8576188,-0.0015 a 1.1814304,1.171762 0 0 0 0.9236308,-0.44011 l 4.2745966,-5.317985 a 1.1814304,1.171762 0 0 0 0.228288,-0.990993 L 0.53894439,7.6775738 A 1.1814304,1.171762 0 0 0 -0.10026101,6.8834313 L -6.2790037,3.9321555 A 1.1814304,1.171762 0 0 0 -6.8523435,3.8176372 Z m 0.00299,0.4550789 a 1.1191255,1.1099671 0 0 1 0.5426517,0.1085303 l 5.85315169,2.7948425 A 1.1191255,1.1099671 0 0 1 0.15197811,7.9290648 L 1.598051,14.21035 a 1.1191255,1.1099671 0 0 1 -0.2163123,0.939348 l -4.0493032,5.037304 a 1.1191255,1.1099671 0 0 1 -0.8749789,0.416906 l -6.4961006,0.0015 a 1.1191255,1.1099671 0 0 1 -0.874979,-0.417652 l -4.0508,-5.036554 a 1.1191255,1.1099671 0 0 1 -0.06362,-0.08832 1.1191255,1.1099671 0 0 1 -0.151942,-0.851028 l 1.443827,-6.2812853 a 1.1191255,1.1099671 0 0 1 0.605524,-0.7529758 l 5.8524036,-2.7963395 a 1.1191255,1.1099671 0 0 1 0.4288819,-0.1085303 z" style="color:#000000;font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:medium;line-height:normal;font-family:Sans;-inkscape-font-specification:Sans;text-indent:0;text-align:start;text-decoration:none;text-decoration-line:none;letter-spacing:normal;word-spacing:normal;text-transform:none;writing-mode:lr-tb;direction:ltr;baseline-shift:baseline;text-anchor:start;display:inline;overflow:visible;visibility:visible;fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;marker:none;enable-background:accumulate" />
</g>
<text y="16.811775" x="9.9744644" style="font-style:normal;font-weight:normal;font-size:10.58333302px;line-height:6.61458349px;font-family:Sans;letter-spacing:0px;word-spacing:0px;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1">
	<tspan y="16.811775" x="9.9744644" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:2.82222223px;font-family:Arial;-inkscape-font-specification:'Arial, Normal';text-align:center;writing-mode:lr-tb;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.26458332px">deploy</tspan>
</text>
<g transform="translate(-0.65385546,0)">
	<path style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" d="m 10.225062,13.731632 0,0 C 7.7824218,13.847177 5.7050116,11.968386 5.5753417,9.5264634 5.4456516,7.0845405 7.3124018,4.9962905 9.7535318,4.8524795 c 2.4411202,-0.143811 4.5401412,1.71081 4.6980812,4.1510682 l -1.757081,0.1137208 c -0.0954,-1.473818 -1.36311,-2.593935 -2.8374602,-2.50708 -1.47434,0.08686 -2.60178,1.3480761 -2.52346,2.8228991 0.0783,1.4748224 1.333,2.6095384 2.8082502,2.5397534 z"/>
	<path style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" d="m 11.135574,9.0088015 1.39745,3.4205085 3.2263,-3.4205085 z" />
</g>
</g>
`)

		defsEl.append("symbol").attr("id", "ac_ic_daemonset")
			.attr("width", "18.035334mm").attr("height", "17.500378mm").attr("viewBox", "0 0 18.035334 17.500378")
			.html(`<g transform="translate(-0.99262638,-1.174181)">
<g transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">
	<path d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z" style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="M -6.8523435,3.8176372 A 1.1814304,1.171762 0 0 0 -7.3044284,3.932904 l -6.1787426,2.9512758 a 1.1814304,1.171762 0 0 0 -0.639206,0.794891 l -1.523915,6.6308282 a 1.1814304,1.171762 0 0 0 0.160175,0.89893 1.1814304,1.171762 0 0 0 0.06736,0.09281 l 4.276094,5.317236 a 1.1814304,1.171762 0 0 0 0.92363,0.440858 l 6.8576188,-0.0015 a 1.1814304,1.171762 0 0 0 0.9236308,-0.44011 l 4.2745966,-5.317985 a 1.1814304,1.171762 0 0 0 0.228288,-0.990993 L 0.53894439,7.6775738 A 1.1814304,1.171762 0 0 0 -0.10026101,6.8834313 L -6.2790037,3.9321555 A 1.1814304,1.171762 0 0 0 -6.8523435,3.8176372 Z m 0.00299,0.4550789 a 1.1191255,1.1099671 0 0 1 0.5426517,0.1085303 l 5.85315169,2.7948425 A 1.1191255,1.1099671 0 0 1 0.15197811,7.9290648 L 1.598051,14.21035 a 1.1191255,1.1099671 0 0 1 -0.2163123,0.939348 l -4.0493032,5.037304 a 1.1191255,1.1099671 0 0 1 -0.8749789,0.416906 l -6.4961006,0.0015 a 1.1191255,1.1099671 0 0 1 -0.874979,-0.417652 l -4.0508,-5.036554 a 1.1191255,1.1099671 0 0 1 -0.06362,-0.08832 1.1191255,1.1099671 0 0 1 -0.151942,-0.851028 l 1.443827,-6.2812853 a 1.1191255,1.1099671 0 0 1 0.605524,-0.7529758 l 5.8524036,-2.7963395 a 1.1191255,1.1099671 0 0 1 0.4288819,-0.1085303 z" style="color:#000000;font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:medium;line-height:normal;font-family:Sans;-inkscape-font-specification:Sans;text-indent:0;text-align:start;text-decoration:none;text-decoration-line:none;letter-spacing:normal;word-spacing:normal;text-transform:none;writing-mode:lr-tb;direction:ltr;baseline-shift:baseline;text-anchor:start;display:inline;overflow:visible;visibility:visible;fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;marker:none;enable-background:accumulate" />
</g>
<text y="16.811775" x="10.016495" style="font-style:normal;font-weight:normal;font-size:10.58333302px;line-height:6.61458349px;font-family:Sans;letter-spacing:0px;word-spacing:0px;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1">
	<tspan y="16.811775" x="10.016495" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:2.82222223px;font-family:Arial;-inkscape-font-specification:'Arial, Normal';text-align:center;writing-mode:lr-tb;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.26458332px">ds</tspan>
</text>
<g transform="translate(0.58627835,0)">
	<path d="m 7.708299,5.2827748 6.524989,0 0,4.5833348 -6.524989,0 z" style="fill:none;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52914584;stroke-linecap:square;stroke-linejoin:round;stroke-miterlimit:10;stroke-dasharray:1.58743756, 1.58743756;stroke-dashoffset:3.66698074;stroke-opacity:1" />
	<path d="m 4.350169,13.606752 7.074559,0" style="fill:none;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.61833036;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="m 6.169549,6.6940855 6.524989,0 0,4.5833355 -6.524989,0 z" style="fill:#326ce5;fill-opacity:1;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52914584;stroke-linecap:square;stroke-linejoin:round;stroke-miterlimit:10;stroke-dasharray:1.58743756, 1.58743756;stroke-dashoffset:3.87863898;stroke-opacity:1" />
	<path d="m 4.630799,8.1053983 6.524999,0 0,4.5833347 -6.524999,0 z" style="fill:none;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52916664;stroke-linecap:butt;stroke-linejoin:round;stroke-miterlimit:10;stroke-opacity:1" />
	<path d="m 4.5865192,8.1226661 6.5250018,0 0,4.5833339 -6.5250018,0 z" style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" />
</g>
</g>`)

		defsEl.append("symbol").attr("id", "ac_ic_replicaset")
			.attr("width", "18.035334mm").attr("height", "17.500378mm").attr("viewBox", "0 0 18.035334 17.500378")
			.html(`<g transform="translate(-0.99262638,-1.174181)">
<g transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">
	<path d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z" style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="M -6.8523435,3.8176372 A 1.1814304,1.171762 0 0 0 -7.3044284,3.932904 l -6.1787426,2.9512758 a 1.1814304,1.171762 0 0 0 -0.639206,0.794891 l -1.523915,6.6308282 a 1.1814304,1.171762 0 0 0 0.160175,0.89893 1.1814304,1.171762 0 0 0 0.06736,0.09281 l 4.276094,5.317236 a 1.1814304,1.171762 0 0 0 0.92363,0.440858 l 6.8576188,-0.0015 a 1.1814304,1.171762 0 0 0 0.9236308,-0.44011 l 4.2745966,-5.317985 a 1.1814304,1.171762 0 0 0 0.228288,-0.990993 L 0.53894439,7.6775738 A 1.1814304,1.171762 0 0 0 -0.10026101,6.8834313 L -6.2790037,3.9321555 A 1.1814304,1.171762 0 0 0 -6.8523435,3.8176372 Z m 0.00299,0.4550789 a 1.1191255,1.1099671 0 0 1 0.5426517,0.1085303 l 5.85315169,2.7948425 A 1.1191255,1.1099671 0 0 1 0.15197811,7.9290648 L 1.598051,14.21035 a 1.1191255,1.1099671 0 0 1 -0.2163123,0.939348 l -4.0493032,5.037304 a 1.1191255,1.1099671 0 0 1 -0.8749789,0.416906 l -6.4961006,0.0015 a 1.1191255,1.1099671 0 0 1 -0.874979,-0.417652 l -4.0508,-5.036554 a 1.1191255,1.1099671 0 0 1 -0.06362,-0.08832 1.1191255,1.1099671 0 0 1 -0.151942,-0.851028 l 1.443827,-6.2812853 a 1.1191255,1.1099671 0 0 1 0.605524,-0.7529758 l 5.8524036,-2.7963395 a 1.1191255,1.1099671 0 0 1 0.4288819,-0.1085303 z" style="color:#000000;font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:medium;line-height:normal;font-family:Sans;-inkscape-font-specification:Sans;text-indent:0;text-align:start;text-decoration:none;text-decoration-line:none;letter-spacing:normal;word-spacing:normal;text-transform:none;writing-mode:lr-tb;direction:ltr;baseline-shift:baseline;text-anchor:start;display:inline;overflow:visible;visibility:visible;fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;marker:none;enable-background:accumulate" />
</g>
<text x="9.9730864" y="16.811775" style="font-style:normal;font-weight:normal;font-size:10.58333302px;line-height:6.61458349px;font-family:Sans;letter-spacing:0px;word-spacing:0px;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1">
	<tspan x="9.9730864" y="16.811775" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:2.82222223px;font-family:Arial;-inkscape-font-specification:'Arial, Normal';text-align:center;writing-mode:lr-tb;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.26458332px">rs</tspan>
</text>
<g transform="translate(0.16298107,0)">
	<path d="m 8.123609,5.5524084 6.52499,0 0,4.5833346 -6.52499,0 z" style="fill:#326ce5;fill-opacity:1;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52899998;stroke-linecap:square;stroke-linejoin:round;stroke-miterlimit:10;stroke-dasharray:1.58700001, 1.58700001;stroke-dashoffset:3.66597009;stroke-opacity:1" />
	<path d="m 6.5848588,6.9637194 6.5249902,0 0,4.5833346 -6.5249902,0 z" style="fill:#326ce5;fill-opacity:1;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52914584;stroke-linecap:square;stroke-linejoin:round;stroke-miterlimit:10;stroke-dasharray:1.58743756, 1.58743756;stroke-dashoffset:3.87863898;stroke-opacity:1" />
	<path d="m 5.0461088,8.3750314 6.5250002,0 0,4.5833346 -6.5250002,0 z" style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" />
	<path d="m 5.0461088,8.3750314 6.5250002,0 0,4.5833346 -6.5250002,0 z" style="fill:none;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52916664;stroke-linecap:butt;stroke-linejoin:round;stroke-miterlimit:10;stroke-opacity:1" />
</g>
</g>`)


		defsEl.append("symbol").attr("id", "ac_ic_pod")
			.attr("width", "18.035334mm").attr("height", "17.500378mm").attr("viewBox", "0 0 18.035334 17.500378")
			.html(`<g transform="translate(-0.99262638,-1.174181)">
<g transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">
	<path d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z" style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="M -6.8523435,3.8176372 A 1.1814304,1.171762 0 0 0 -7.3044284,3.932904 l -6.1787426,2.9512758 a 1.1814304,1.171762 0 0 0 -0.639206,0.794891 l -1.523915,6.6308282 a 1.1814304,1.171762 0 0 0 0.160175,0.89893 1.1814304,1.171762 0 0 0 0.06736,0.09281 l 4.276094,5.317236 a 1.1814304,1.171762 0 0 0 0.92363,0.440858 l 6.8576188,-0.0015 a 1.1814304,1.171762 0 0 0 0.9236308,-0.44011 l 4.2745966,-5.317985 a 1.1814304,1.171762 0 0 0 0.228288,-0.990993 L 0.53894439,7.6775738 A 1.1814304,1.171762 0 0 0 -0.10026101,6.8834313 L -6.2790037,3.9321555 A 1.1814304,1.171762 0 0 0 -6.8523435,3.8176372 Z m 0.00299,0.4550789 a 1.1191255,1.1099671 0 0 1 0.5426517,0.1085303 l 5.85315169,2.7948425 A 1.1191255,1.1099671 0 0 1 0.15197811,7.9290648 L 1.598051,14.21035 a 1.1191255,1.1099671 0 0 1 -0.2163123,0.939348 l -4.0493032,5.037304 a 1.1191255,1.1099671 0 0 1 -0.8749789,0.416906 l -6.4961006,0.0015 a 1.1191255,1.1099671 0 0 1 -0.874979,-0.417652 l -4.0508,-5.036554 a 1.1191255,1.1099671 0 0 1 -0.06362,-0.08832 1.1191255,1.1099671 0 0 1 -0.151942,-0.851028 l 1.443827,-6.2812853 a 1.1191255,1.1099671 0 0 1 0.605524,-0.7529758 l 5.8524036,-2.7963395 a 1.1191255,1.1099671 0 0 1 0.4288819,-0.1085303 z" style="color:#000000;font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:medium;line-height:normal;font-family:Sans;-inkscape-font-specification:Sans;text-indent:0;text-align:start;text-decoration:none;text-decoration-line:none;letter-spacing:normal;word-spacing:normal;text-transform:none;writing-mode:lr-tb;direction:ltr;baseline-shift:baseline;text-anchor:start;display:inline;overflow:visible;visibility:visible;fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;marker:none;enable-background:accumulate" />
</g>
<text x="10.017183" y="16.811775" style="font-style:normal;font-weight:normal;font-size:10.58333302px;line-height:6.61458349px;font-family:Sans;letter-spacing:0px;word-spacing:0px;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1">
	<tspan x="10.017183" y="16.811775" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:2.82222223px;font-family:Arial;-inkscape-font-specification:'Arial, Normal';text-align:center;writing-mode:lr-tb;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.26458332px">pod</tspan>
</text>
<g transform="translate(0.12766661,0)">
	<path d="M 6.2617914,7.036086 9.8826317,5.986087 13.503462,7.036086 9.8826317,8.086087 Z" style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" />
	<path d="m 6.2617914,7.43817 0,3.852778 3.3736103,1.868749 0.0167,-4.713193 z" style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" />
	<path d="m 13.503462,7.43817 0,3.852778 -3.37361,1.868749 -0.0167,-4.713193 z" style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" />
</g>
</g>`)
		defsEl.append("symbol").attr("id", "ac_ic_ingress")
			.attr("width", "18.035334mm").attr("height", "17.500378mm").attr("viewBox", "0 0 18.035334 17.500378")
			.html(`<g transform="translate(-0.99262638,-1.174181)">
<g	transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">
	<path d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z" style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="M -6.8523435,3.8176372 A 1.1814304,1.171762 0 0 0 -7.3044284,3.932904 l -6.1787426,2.9512758 a 1.1814304,1.171762 0 0 0 -0.639206,0.794891 l -1.523915,6.6308282 a 1.1814304,1.171762 0 0 0 0.160175,0.89893 1.1814304,1.171762 0 0 0 0.06736,0.09281 l 4.276094,5.317236 a 1.1814304,1.171762 0 0 0 0.92363,0.440858 l 6.8576188,-0.0015 a 1.1814304,1.171762 0 0 0 0.9236308,-0.44011 l 4.2745966,-5.317985 a 1.1814304,1.171762 0 0 0 0.228288,-0.990993 L 0.53894439,7.6775738 A 1.1814304,1.171762 0 0 0 -0.10026101,6.8834313 L -6.2790037,3.9321555 A 1.1814304,1.171762 0 0 0 -6.8523435,3.8176372 Z m 0.00299,0.4550789 a 1.1191255,1.1099671 0 0 1 0.5426517,0.1085303 l 5.85315169,2.7948425 A 1.1191255,1.1099671 0 0 1 0.15197811,7.9290648 L 1.598051,14.21035 a 1.1191255,1.1099671 0 0 1 -0.2163123,0.939348 l -4.0493032,5.037304 a 1.1191255,1.1099671 0 0 1 -0.8749789,0.416906 l -6.4961006,0.0015 a 1.1191255,1.1099671 0 0 1 -0.874979,-0.417652 l -4.0508,-5.036554 a 1.1191255,1.1099671 0 0 1 -0.06362,-0.08832 1.1191255,1.1099671 0 0 1 -0.151942,-0.851028 l 1.443827,-6.2812853 a 1.1191255,1.1099671 0 0 1 0.605524,-0.7529758 l 5.8524036,-2.7963395 a 1.1191255,1.1099671 0 0 1 0.4288819,-0.1085303 z" style="color:#000000;font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:medium;line-height:normal;font-family:Sans;-inkscape-font-specification:Sans;text-indent:0;text-align:start;text-decoration:none;text-decoration-line:none;letter-spacing:normal;word-spacing:normal;text-transform:none;writing-mode:lr-tb;direction:ltr;baseline-shift:baseline;text-anchor:start;display:inline;overflow:visible;visibility:visible;fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;marker:none;enable-background:accumulate"/>
</g>
<text x="10.008915" y="16.811775" style="font-style:normal;font-weight:normal;font-size:10.58333302px;line-height:6.61458349px;font-family:Sans;letter-spacing:0px;word-spacing:0px;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1">
	<tspan x="10.008915" y="16.811775" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:2.82222223px;font-family:Arial;-inkscape-font-specification:'Arial, Normal';text-align:center;writing-mode:lr-tb;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.26458332px">ing</tspan>
</text>
<path d="m 12.75799,13.716256 -2.270701,0 -4.9209009,-6.1558614 -1.42366,0 0,-2.0149069 2.31473,0 4.9230119,6.1558533 1.37752,0 0,-1.593474 3.119869,2.599882 -3.119869,2.601983 z m -2.47616,-4.7552748 1.09864,-1.3754256 1.37752,0 0,1.593475 3.119869,-2.5998829 -3.119869,-2.601983 0,1.593483 -2.270701,0 -1.4571904,1.8241102 z m -3.5979219,1.3649428 -1.11752,1.400578 -1.42366,0 0,2.014915 2.31473,0 1.4781699,-1.849278 z"style="fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.20966817" />
</g>`)

		defsEl.append("symbol").attr("id", "ac_ic_service")
			.attr("width", "18.035334mm").attr("height", "17.500378mm").attr("viewBox", "0 0 18.035334 17.500378")
			.html(`<g transform="translate(-0.99262638,-1.174181)">
<g transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">
	<path d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z" style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />
	<path d="M -6.8523435,3.8176372 A 1.1814304,1.171762 0 0 0 -7.3044284,3.932904 l -6.1787426,2.9512758 a 1.1814304,1.171762 0 0 0 -0.639206,0.794891 l -1.523915,6.6308282 a 1.1814304,1.171762 0 0 0 0.160175,0.89893 1.1814304,1.171762 0 0 0 0.06736,0.09281 l 4.276094,5.317236 a 1.1814304,1.171762 0 0 0 0.92363,0.440858 l 6.8576188,-0.0015 a 1.1814304,1.171762 0 0 0 0.9236308,-0.44011 l 4.2745966,-5.317985 a 1.1814304,1.171762 0 0 0 0.228288,-0.990993 L 0.53894439,7.6775738 A 1.1814304,1.171762 0 0 0 -0.10026101,6.8834313 L -6.2790037,3.9321555 A 1.1814304,1.171762 0 0 0 -6.8523435,3.8176372 Z m 0.00299,0.4550789 a 1.1191255,1.1099671 0 0 1 0.5426517,0.1085303 l 5.85315169,2.7948425 A 1.1191255,1.1099671 0 0 1 0.15197811,7.9290648 L 1.598051,14.21035 a 1.1191255,1.1099671 0 0 1 -0.2163123,0.939348 l -4.0493032,5.037304 a 1.1191255,1.1099671 0 0 1 -0.8749789,0.416906 l -6.4961006,0.0015 a 1.1191255,1.1099671 0 0 1 -0.874979,-0.417652 l -4.0508,-5.036554 a 1.1191255,1.1099671 0 0 1 -0.06362,-0.08832 1.1191255,1.1099671 0 0 1 -0.151942,-0.851028 l 1.443827,-6.2812853 a 1.1191255,1.1099671 0 0 1 0.605524,-0.7529758 l 5.8524036,-2.7963395 a 1.1191255,1.1099671 0 0 1 0.4288819,-0.1085303 z" style="color:#000000;font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:medium;line-height:normal;font-family:Sans;-inkscape-font-specification:Sans;text-indent:0;text-align:start;text-decoration:none;text-decoration-line:none;letter-spacing:normal;word-spacing:normal;text-transform:none;writing-mode:lr-tb;direction:ltr;baseline-shift:baseline;text-anchor:start;display:inline;overflow:visible;visibility:visible;fill:#ffffff;fill-opacity:1;fill-rule:nonzero;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;marker:none;enable-background:accumulate" />
</g>
<text x="10.008915" y="16.811775" style="font-style:normal;font-weight:normal;font-size:10.58333302px;line-height:6.61458349px;font-family:Sans;letter-spacing:0px;word-spacing:0px;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:0.26458332px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1">
	<tspan x="9.976531" y="16.811775" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:2.82222223px;font-family:Arial;-inkscape-font-specification:'Arial, Normal';text-align:center;writing-mode:lr-tb;text-anchor:middle;fill:#ffffff;fill-opacity:1;stroke-width:0.26458332px">svc</tspan>
</text>
<g transform="translate(0.09238801,0)">
	<path style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" d="m 4.4949896,11.260826 2.9083311,0 0,2.041667 -2.9083311,0 z"/>
	<path style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" d="m 8.4637407,11.260826 2.9083303,0 0,2.041667 -2.9083303,0 z"/>
	<path style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" d="m 12.432491,11.260826 2.90833,0 0,2.041667 -2.90833,0 z"/>
	<path style="fill:#ffffff;fill-rule:evenodd;stroke:none;stroke-width:0.26458332;stroke-linecap:square;stroke-miterlimit:10" d="m 7.6137407,5.2082921 4.6083303,0 0,2.041667 -4.6083303,0 z"/>
	<path style="fill:none;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52916664;stroke-linecap:butt;stroke-linejoin:round;stroke-miterlimit:10;stroke-opacity:1" d="m 9.9179005,7.2499601 0,2.005449 -3.966671,0 0,2.0028859"/>
	<path style="fill:none;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52899998;stroke-linecap:butt;stroke-linejoin:round;stroke-miterlimit:10;stroke-dasharray:none;stroke-opacity:1" d="m 9.9179005,7.2499601 0,2.005449 3.9666705,0 0,2.0028859"/>
	<path style="fill:none;fill-rule:evenodd;stroke:#ffffff;stroke-width:0.52916664;stroke-linecap:butt;stroke-linejoin:round;stroke-miterlimit:10;stroke-opacity:1" d="m 9.9095538,7.2512251 0,2.005449 0.0167,0 0,2.0028859"/>
</g>
</g>`)


	}

};	
