import{_ as i,o as a,c as n,ag as p}from"./chunks/framework.DXGyWiRo.js";const t="/images/14.6.pprof.png?raw=true",l="/images/14.6.pprof2.png?raw=true",e="/images/14.6.pprof3.png?raw=true",y=JSON.parse('{"title":"14.6 pprof支持","description":"","frontmatter":{},"headers":[],"relativePath":"14/6.md","filePath":"14/6.md","lastUpdated":null}'),h={name:"14/6.md"};function k(r,s,o,E,d,g){return a(),n("div",null,[...s[0]||(s[0]=[p(`<h1 id="_14-6-pprof支持" tabindex="-1">14.6 pprof支持 <a class="header-anchor" href="#_14-6-pprof支持" aria-label="Permalink to &quot;14.6 pprof支持&quot;">​</a></h1><p>Go语言有一个非常棒的设计就是标准库里面带有代码的性能监控工具，在两个地方有包：</p><div class="language-Go vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">Go</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">net</span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">/</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">http</span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">/</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">pprof</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">runtime</span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">/</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">pprof</span></span></code></pre></div><p>其实net/http/pprof中只是使用runtime/pprof包来进行封装了一下，并在http端口上暴露出来</p><h2 id="beego支持pprof" tabindex="-1">beego支持pprof <a class="header-anchor" href="#beego支持pprof" aria-label="Permalink to &quot;beego支持pprof&quot;">​</a></h2><p>目前beego框架新增了pprof，该特性默认是不开启的，如果你需要测试性能，查看相应的执行goroutine之类的信息，其实Go的默认包&quot;net/http/pprof&quot;已经具有该功能，如果按照Go默认的方式执行Web，默认就可以使用，但是由于beego重新封装了ServHTTP函数，默认的包是无法开启该功能的，所以需要对beego的内部改造支持pprof。</p><ul><li>首先在beego.Run函数中根据变量是否自动加载性能包</li></ul><div class="language-Go vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">Go</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">if</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> PprofOn {</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	BeeApp.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">RegisterController</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">\`/debug/pprof\`</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">, </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">&amp;</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">ProfController</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">{})</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	BeeApp.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">RegisterController</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">\`/debug/pprof/:pp([\\w]+)\`</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">, </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">&amp;</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">ProfController</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">{})</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">}</span></span></code></pre></div><ul><li>设计ProfController</li></ul><div class="language-Go vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">Go</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">package</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;"> beego</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">import</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> (</span></span>
<span class="line"><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">	&quot;</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">net/http/pprof</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">)</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">type</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;"> ProfController</span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;"> struct</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> {</span></span>
<span class="line"><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">	Controller</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">}</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">func</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> (</span><span style="--shiki-light:#E36209;--shiki-dark:#FFAB70;">this </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">*</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">ProfController</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">) </span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Get</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">() {</span></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">	switch</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> this.Ctx.Param[</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">&quot;:pp&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">] {</span></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">	default</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">:</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">		pprof.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Index</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(this.Ctx.ResponseWriter, this.Ctx.Request)</span></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">	case</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> &quot;&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">:</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">		pprof.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Index</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(this.Ctx.ResponseWriter, this.Ctx.Request)</span></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">	case</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> &quot;cmdline&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">:</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">		pprof.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Cmdline</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(this.Ctx.ResponseWriter, this.Ctx.Request)</span></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">	case</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> &quot;profile&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">:</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">		pprof.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Profile</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(this.Ctx.ResponseWriter, this.Ctx.Request)</span></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">	case</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> &quot;symbol&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">:</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">		pprof.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Symbol</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(this.Ctx.ResponseWriter, this.Ctx.Request)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	}</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	this.Ctx.ResponseWriter.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">WriteHeader</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">(</span><span style="--shiki-light:#005CC5;--shiki-dark:#79B8FF;">200</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">}</span></span></code></pre></div><h2 id="使用入门" tabindex="-1">使用入门 <a class="header-anchor" href="#使用入门" aria-label="Permalink to &quot;使用入门&quot;">​</a></h2><p>通过上面的设计，你可以通过如下代码开启pprof：</p><div class="language-Go vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">Go</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">beego.PprofOn </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">=</span><span style="--shiki-light:#005CC5;--shiki-dark:#79B8FF;"> true</span></span></code></pre></div><p>然后你就可以在浏览器中打开如下URL就看到如下界面： <img src="`+t+'" alt=""></p><p>图14.7 系统当前goroutine、heap、thread信息</p><p>点击goroutine我们可以看到很多详细的信息：</p><p><img src="'+l+`" alt=""></p><p>图14.8 显示当前goroutine的详细信息</p><p>我们还可以通过命令行获取更多详细的信息</p><div class="language-Go vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">Go</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">go</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> tool pprof http:</span><span style="--shiki-light:#6A737D;--shiki-dark:#6A737D;">//localhost:8080/debug/pprof/profile</span></span></code></pre></div><p>这时候程序就会进入30秒的profile收集时间，在这段时间内拼命刷新浏览器上的页面，尽量让cpu占用性能产生数据。</p><pre><code>(pprof) top10

Total: 3 samples

   1 33.3% 33.3% 1 33.3% MHeap_AllocLocked

   1 33.3% 66.7% 1 33.3% os/exec.(*Cmd).closeDescriptors

   1 33.3% 100.0% 1 33.3% runtime.sigprocmask

   0 0.0% 100.0% 1 33.3% MCentral_Grow

   0 0.0% 100.0% 2 66.7% main.Compile

   0 0.0% 100.0% 2 66.7% main.compile

   0 0.0% 100.0% 2 66.7% main.run

   0 0.0% 100.0% 1 33.3% makeslice1

   0 0.0% 100.0% 2 66.7% net/http.(*ServeMux).ServeHTTP

   0 0.0% 100.0% 2 66.7% net/http.(*conn).serve	

(pprof)web
</code></pre><p><img src="`+e+'" alt=""></p><p>图14.9 展示的执行流程信息</p><h2 id="links" tabindex="-1">links <a class="header-anchor" href="#links" aria-label="Permalink to &quot;links&quot;">​</a></h2><ul><li><a href="/">目录</a></li><li>上一节: <a href="/14/5.html">多语言支持</a></li><li>下一节: <a href="/14/7.html">小结</a></li></ul>',26)])])}const u=i(h,[["render",k]]);export{y as __pageData,u as default};
