import{r,j as e,P as $,i as L}from"./collect-ui-core-CiE9scW-.js";import{F as G,a as J,b as Q}from"./icons-CBicxOix.js";import{g as D,a as H}from"./getFileName-CES2ZA2y.js";import{r as X,u as Y,D as Z}from"./office-Bb9dSmE3.js";import{S as O,A as ee,R as N,B as A,a as te,T as ae,b as W,c as ne,d as se,U as re,I as oe,e as le,f as ie,g as ce,h as de,i as fe,j as xe,k as he,l as ue,m as ge,n as me}from"./antd-koSoMy5T.js";const ye=({url:t})=>{const[w,v]=r.useState([]),[j,l]=r.useState(""),[h,b]=r.useState({}),[T,I]=r.useState(!0),[i,C]=r.useState(null),u=(s,c)=>s?`
      <div class="excel-container">
        ${s}
      </div>
      <style>
        /* 容器样式 */
        .excel-container {
          overflow-x: auto;
          max-height: 70vh;
          border: 1px solid #e0e0e0;
          background: white;
          font-family: Arial;
        }

        /* 表格样式 */
        .excel-container table {
          border-collapse: collapse;
          min-width: max-content !important;
          width: 100%;
        }

        /* 表头固定 */
        .excel-container th {
          position: sticky;
          top: 0;
          background: #f0f0f0 !important;
          z-index: 2;
          font-weight: bold;
          box-shadow: 0 1px 0 #d9d9d9;
        }

        /* 单元格样式 */
        .excel-container th,
        .excel-container td {
          border: 1px solid #e0e0e0 !important;
          padding: 8px 12px;
          white-space: nowrap;
        }

        /* 斑马纹 */
        .excel-container tr:nth-child(even) {
          background: #f9f9f9 !important;
        }

        /* 悬停效果 */
        .excel-container tr:hover {
          background: #f0f0f0 !important;
        }
         .excel-preview .ant-tabs-nav{
            margin-bottom: 0px;
        
        }
        /* 隐藏首行空单元格（兼容性处理） */
.excel-container tr:first-child td:empty {
  display: none;
}
      </style>
    `:`
        <div class="excel-error">
          <div class="error-message">
            <h4><ExclamationCircleOutlined /> 工作表 "${c}" 无有效数据</h4>
            <p>可能原因：隐藏工作表/图表工作表/空工作表</p>
          </div>
        </div>
        <style>
          .excel-error {
            padding: 24px;
            text-align: center;
            color: #f5222d;
          }
          .excel-error svg {
            margin-right: 8px;
          }
        </style>
      `;return r.useEffect(()=>{(async()=>{try{if(I(!0),C(null),!t||typeof t!="string")throw new Error("无效的文件URL");const c=await fetch(t);if(!c.ok)throw new Error(`请求失败: ${c.status}`);const g=c.headers.get("content-type")||"";if(!(/excel|spreadsheet/.test(g)||t.toLowerCase().endsWith(".xlsx")||t.toLowerCase().endsWith(".xls")))throw new Error("不是有效的Excel文件");const _=await c.arrayBuffer();let o;try{o=X(_,{type:"array"})}catch(d){throw new Error(`解析失败: ${d.message}`)}if(!o.SheetNames?.length)throw new Error("Excel文件中未找到工作表");const x={};o.SheetNames.forEach(d=>{const k=o.Sheets[d];if(!k||!k["!ref"]){x[d]=u(null,d);return}try{let m=Y.sheet_to_html(k,{raw:!0,header:!1,display:!1});m=(y=>{const E=y.indexOf("false"),z=y.indexOf("<table");return E!==-1&&E<z?y.substring(0,E)+y.substring(E+5):y})(m),x[d]=u(m,d)}catch(m){const F=m.message;x[d]=`
              <div class="excel-error">
                <div class="error-message">
                  <h4><ExclamationCircleOutlined /> 工作表 "${d}" 渲染失败</h4>
                  <p>错误详情: ${F}</p>
                  <p>建议: 请在Excel中检查此工作表内容</p>
                </div>
              </div>
            `}}),v(o.SheetNames),l(o.SheetNames[0]),b(x)}catch(c){console.error("Excel加载错误:",c),C(c.message||"未知错误")}finally{I(!1)}})()},[t]),T?e.jsx("div",{style:{textAlign:"center",padding:"40px 0"},children:e.jsx(O,{tip:"正在加载Excel文件...",size:"large"})}):i?e.jsx(ee,{type:"error",message:"Excel文件加载失败",description:e.jsxs("div",{style:{marginTop:16},children:[e.jsxs("p",{children:[e.jsx(N,{})," 错误信息: ",i]}),e.jsx(A,{type:"primary",icon:e.jsx(te,{}),onClick:()=>window.open(t,"_blank"),style:{marginTop:8},children:"下载原始文件检查"})]}),showIcon:!0}):e.jsx("div",{style:{background:"#fff",padding:0,borderRadius:4,boxShadow:"0 1px 3px rgba(0,0,0,0.1)"},children:e.jsx(ae,{activeKey:j,onChange:l,size:"small",type:"card",className:"excel-preview",items:w.map(s=>({key:s,label:e.jsxs("span",{children:[e.jsx(W,{style:{marginRight:8}}),s,h[s]?.includes("excel-error")&&e.jsx(N,{style:{color:"#f5222d",marginLeft:8}})]}),children:e.jsx("div",{dangerouslySetInnerHTML:{__html:h[s]||u(null,s)},style:{marginTop:16}})}))})})},pe=({url:t,style:w})=>{const[v,j]=r.useState(""),[l,h]=r.useState(!0),[b,T]=r.useState(null);return r.useEffect(()=>{(async()=>{try{h(!0);const i=await fetch(t);if(i.ok){const C=await i.text();j(C)}else throw new Error(`Failed to load text file: ${i.status}`)}catch(i){T(i instanceof Error?i.message:"Failed to load text file")}finally{h(!1)}})()},[t]),l?e.jsx("div",{style:{...w,display:"flex",justifyContent:"center",alignItems:"center"},children:e.jsx(O,{tip:"Loading text content..."})}):b?e.jsxs("div",{style:{...w,color:"red",padding:16},children:["Error: ",b]}):e.jsx("div",{style:{...w,overflow:"auto",padding:16,backgroundColor:"#f5f5f5",whiteSpace:"pre-wrap",fontFamily:"monospace",border:"1px solid #d9d9d9",borderRadius:4,lineHeight:1.5},children:v||"Empty text file"})};function ke(t){const{attachment_prop:w,show_path:v,finish_action:j,uploadConfig:l,placeholder:h,...b}=t,{visible:T,...I}=$.transferProp(b,"attachment"),i=L("dialog");L("icon");const[C,u]=r.useState(!1),[s,c]=r.useState(""),g=ne.useApp(),R=$.toApiObj(l?.api);let _={};if(R?.data)for(let a in R.data)_[a]=R.data[a];if(l?.data)for(let a in l?.data){const n=l?.data[a];_[a]=$.varValue(n,t.store,I.target)}const o=l?.multiple||!1,x=t?._target?.row[t?._target?.column?.field],d=()=>{if(!s)return null;switch(H(D(s))){case"word":return e.jsx(Z,{url:s,style:{height:"80vh",overflow:"auto"}});case"pdf":return e.jsx("iframe",{src:s,style:{height:"80vh",width:"100%"}});case"excel":return e.jsx(ye,{url:s});case"properties":case"json":case"xml":case"sql":case"yml":case"text":return e.jsx(pe,{url:s,style:{height:"80vh"}});default:return null}},k=a=>{if(!a)return!1;const n=H(D(a));return["word","pdf","excel","text","properties","json","xml","sql","yml"].includes(n)},m=r.useCallback(a=>{if(a.file.status==="done"){console.log("上传成功后的返回数据:",a.file.response),j&&$.handlerActions(j,t.store,t.rootStore,g,{row:a.file.response});const n=a.file.response.data;t.onChange&&(o?t.onChange([n,...t?.value||[]]):t.onChange(n?.path)),t?._target?.onValueChange&&(t?._target?.onValueChange(n?.path),t?._target?.api.stopEditing())}else a.file.status==="error"&&g?.message?.error(`${a.file.name} 文件上传失败`)},[t.value,o]),F=r.useCallback(a=>{if(t?._target?.onValueChange){if(o){const n=Array.isArray(t.value)?[...t.value].filter((p,f)=>f!==a):[];t._target.onValueChange(n)}else t._target.onValueChange("");t._target.api.stopEditing()}if(t.onChange)if(o){const n=Array.isArray(t.value)?[...t.value].filter((p,f)=>f!==a):[];t.onChange(n)}else t.onChange("")},[t.value,o,t.onChange,t._target]),y=a=>({word:e.jsx(me,{}),pdf:e.jsx(ge,{}),excel:e.jsx(W,{}),ppt:e.jsx(ue,{}),image:e.jsx(he,{}),video:e.jsx(xe,{}),audio:e.jsx(fe,{}),zip:e.jsx(de,{})})[a]||e.jsx(ce,{}),E=a=>({word:"#2b579a",pdf:"#d24726",excel:"#217346",ppt:"#d24726",zip:"#7e57c2"})[a]||"#999",z=r.useCallback((a,n)=>{if(!a){g?.message?.error("没有可下载的文件");return}try{const p=n||D(a),f=document.createElement("a");f.href=a,f.download=p,document.body.appendChild(f),f.click(),document.body.removeChild(f),g?.message?.success(`开始下载: ${p}`)}catch{g?.message?.error("下载文件时出错")}},[]),U=({file:a,index:n,onPreview:p,onDownload:f,onRemove:M,showPreview:K=!0})=>{const S=a.path||a;debugger;const V=D(S),P=H(V);return e.jsxs("div",{style:{width:120,height:120,border:"1px solid #d9d9d9",borderRadius:4,padding:8,display:"flex",flexDirection:"column",position:"relative",backgroundColor:"#fafafa",overflow:"hidden"},children:[e.jsx("div",{style:{height:80,display:"flex",justifyContent:"center",alignItems:"center",overflow:"hidden"},children:P==="image"?e.jsx("div",{style:{display:"inline-flex",maxWidth:"100%",maxHeight:"100%"},children:e.jsx(le,{src:S,style:{maxWidth:"100%",maxHeight:"100%",objectFit:"contain",display:"block",borderRadius:4,cursor:"pointer"},preview:{mask:null,src:S}})}):e.jsx("div",{style:{fontSize:48,color:E(P),textAlign:"center"},children:y(P)})}),e.jsx("div",{style:{marginTop:8,whiteSpace:"nowrap",overflow:"hidden",textOverflow:"ellipsis",fontSize:12,textAlign:"center"},title:V,children:V}),e.jsxs("div",{style:{position:"absolute",bottom:0,right:4,display:"flex",gap:4},children:[K&&k(S)&&e.jsx(A,{size:"small",type:"text",icon:e.jsx(ie,{}),onClick:()=>p(S),style:{color:"#1890ff"}}),e.jsx(A,{size:"small",type:"text",icon:e.jsx(J,{}),onClick:()=>f(S,V),style:{color:"#52c41a"}}),e.jsx(A,{size:"small",type:"text",icon:e.jsx(Q,{}),onClick:()=>M(n),style:{color:"#ff4d4f"}})]})]})},B=r.useCallback(a=>{c(a),u(!0)},[]),q=r.useCallback(()=>o&&Array.isArray(t.value)?t.value:t.value||x?[t.value||x]:[],[o,t.value,x]);return $.getVisible(t)?e.jsxs(e.Fragment,{children:[e.jsxs(se.Compact,{style:{width:"100%"},children:[e.jsx(re,{...l,onChange:m,name:"file",action:R?.url,data:a=>({..._}),children:e.jsx(A,{icon:e.jsx(G,{}),type:"primary",children:!v&&"上传文件"})}),v&&e.jsx(e.Fragment,{children:e.jsx(oe,{value:t.value||x,onChange:t.onChange||t?._target?.onValueChange,placeholder:h})})]}),e.jsx("div",{style:{marginTop:16,display:"flex",flexWrap:"wrap",gap:12,maxWidth:"100%"},children:q().map((a,n)=>e.jsx(U,{file:a,index:n,onPreview:B,onDownload:z,onRemove:F},n))}),k(s)&&e.jsx(i,{width:"80%",style:{top:"20px"},onOk:()=>u(!1),onCancel:()=>u(!1),open:C,title:"预览文档",children:d()})]}):null}export{ke as default};
