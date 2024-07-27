$(function(){
    baseApp.init();
    $(window).resize(function(){
        baseApp.resizeIframe();
    })
})
var baseApp={
    init:function(){
        this.initAside()
        this.confirmDelete()
        this.resizeIframe()
        this.changeStatus()
        this.changeNum()
    },
    initAside:function(){
        $('.aside h4').click(function(){
            $(this).siblings('ul').slideToggle();
        })
    },
    //设置iframe的高度
    resizeIframe:function(){
        $("#rightMain").height($(window).height()-80)
    },
    // 删除提示
    confirmDelete:function(){
        $(".delete").click(function(){
            var flag=confirm("您确定要删除吗?")
            return flag
        })
    },
    // 更改状态
    changeStatus:function(){
        $(".chStatus").click(function (){
            var id = $(this).attr("data-id");
            var table = $(this).attr("data-table");
            var field = $(this).attr("data-field");
            var el = $(this)
            $.get("/admin/changeStatus", {id:id, table:table, field:field}, function(response){
                console.log(response)
                if(response.success){
                    if (el.attr("src").indexOf("yes")!=-1){
                        el.attr("src","/static/admin/images/no.gif");
                    }else{
                        el.attr("src","/static/admin/images/yes.gif");
                    }
                }
            })
        })
    },
    // 更改排序值
    changeNum:function(){
        $(".chSpanNum").click(function(){
            // 1、获取el以及el里面的属性值
            var id = $(this).attr("data-id");
            var table = $(this).attr("data-table");
            var field = $(this).attr("data-field");
            var num = $(this).html().trim()
            var spanEl = $(this)
            // 2、创建一个  input  的dom节点
            var input = $("<input style='width:60px' value='' />");
            // 3、把  input  放到el里面
            $(this).html(input);
            // 4、让  input  获取焦点,  给  input  赋值
            $(input).trigger("focus").val(num);
            // 5、点击  input  的时候阻止冒泡  （防止触发  changeNum 事件）
            $(input).click(function(e){
                e.stopPropagation();
            })
            // 6、鼠标离开的时候给  span  赋值，并触发 ajax 请求
            $(input).blur(function(){
                var inputNum = $(this).val()
                spanEl.html(inputNum)
                // 触发ajax请求
                $.get("/admin/changeNum", {id:id, table:table, field:field, num:inputNum}, function(response){
                    console.log(response)
                })
            })
        })
    }
}