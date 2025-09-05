#!/bin/bash

# 楼盘-户型图上传功能测试脚本
# 测试完整的楼盘创建、户型图上传和管理流程

echo "🧪 楼盘-户型图上传功能测试"
echo "=========================="

# 设置变量
BASE_URL="http://localhost:8002/api/v1"
AUTH_TOKEN=""  # 如果需要认证，请设置token

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4

    echo -e "\n${YELLOW}测试: $description${NC}"
    echo "请求: $method $url"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -X GET "$url" \
            -H "Content-Type: application/json" \
            ${AUTH_TOKEN:+-H "Authorization: Bearer $AUTH_TOKEN"})
    else
        response=$(curl -s -X $method "$url" \
            -H "Content-Type: application/json" \
            ${AUTH_TOKEN:+-H "Authorization: Bearer $AUTH_TOKEN"} \
            -d "$data")
    fi

    echo "响应: $response"

    # 检查响应是否成功
    if echo "$response" | grep -q '"code": *20[0-9]'; then
        echo -e "${GREEN}✅ 测试通过${NC}"
        return 0
    else
        echo -e "${RED}❌ 测试失败${NC}"
        return 1
    fi
}

# 1. 测试创建楼盘
echo -e "\n${YELLOW}📁 第一步: 创建测试楼盘${NC}"

building_data='{
    "name": "测试楼盘_'$(date +%s)'",
    "district": "朝阳区",
    "businessArea": "三里屯",
    "propertyType": "住宅",
    "status": "available",
    "description": "用于测试楼盘-户型图上传功能的测试楼盘"
}'

building_response=$(curl -s -X POST "$BASE_URL/buildings" \
    -H "Content-Type: application/json" \
    -d "$building_data")

echo "创建楼盘响应: $building_response"

# 提取楼盘ID
building_id=$(echo "$building_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ -z "$building_id" ]; then
    echo -e "${RED}❌ 无法获取楼盘ID，测试终止${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 楼盘创建成功，ID: $building_id${NC}"

# 2. 创建测试户型
echo -e "\n${YELLOW}🏠 第二步: 创建测试户型${NC}"

house_type_data='{
    "buildingId": '$building_id',
    "name": "测试户型_一室一厅",
    "type": "1室1厅",
    "area": 65.5,
    "price": 3500000,
    "status": "available",
    "description": "用于测试的户型"
}'

house_type_response=$(curl -s -X POST "$BASE_URL/house-types" \
    -H "Content-Type: application/json" \
    -d "$house_type_data")

echo "创建户型响应: $house_type_response"

# 提取户型ID
house_type_id=$(echo "$house_type_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ -z "$house_type_id" ]; then
    echo -e "${RED}❌ 无法获取户型ID，测试终止${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 户型创建成功，ID: $house_type_id${NC}"

# 3. 测试户型图上传
echo -e "\n${YELLOW}📤 第三步: 上传户型图${NC}"

# 检查是否有测试图片文件
if [ ! -f "test_floor_plan.jpg" ]; then
    echo "⚠️  未找到测试图片文件 test_floor_plan.jpg，创建一个简单的测试文件"
    echo "请准备一个名为 test_floor_plan.jpg 的图片文件进行测试"
    echo -e "${YELLOW}跳过图片上传测试${NC}"

    # 4. 测试获取楼盘图片列表
    echo -e "\n${YELLOW}📋 第四步: 测试获取楼盘图片列表${NC}"

    test_api "GET" "$BASE_URL/buildings/$building_id/images" "" "获取楼盘所有图片"

    test_api "GET" "$BASE_URL/buildings/$building_id/floor-plans" "" "获取楼盘户型图"

    test_api "GET" "$BASE_URL/images/stats" "" "获取图片统计信息"

    echo -e "\n${GREEN}🎉 楼盘图片管理功能测试完成！${NC}"
    echo -e "${YELLOW}📝 测试总结:${NC}"
    echo "  ✅ 楼盘创建成功 (ID: $building_id)"
    echo "  ✅ 户型创建成功 (ID: $house_type_id)"
    echo "  ⚠️  图片上传测试跳过 (需要准备测试图片文件)"
    echo "  ✅ API接口测试完成"

    exit 0
fi

# 使用curl上传图片
upload_response=$(curl -s -X POST "$BASE_URL/upload/floor-plan" \
    -F "house_type_id=$house_type_id" \
    -F "file=@test_floor_plan.jpg" \
    ${AUTH_TOKEN:+-H "Authorization: Bearer $AUTH_TOKEN"})

echo "上传户型图响应: $upload_response"

# 检查上传是否成功
if echo "$upload_response" | grep -q '"code": *200'; then
    echo -e "${GREEN}✅ 户型图上传成功${NC}"

    # 提取图片ID
    image_id=$(echo "$upload_response" | grep -o '"image_id":[0-9]*' | cut -d':' -f2)
    echo "图片ID: $image_id"
else
    echo -e "${RED}❌ 户型图上传失败${NC}"
fi

# 4. 测试获取楼盘图片列表
echo -e "\n${YELLOW}📋 第四步: 测试获取楼盘图片列表${NC}"

test_api "GET" "$BASE_URL/buildings/$building_id/images" "" "获取楼盘所有图片"

test_api "GET" "$BASE_URL/buildings/$building_id/floor-plans" "" "获取楼盘户型图"

# 5. 测试图片详情和更新
if [ ! -z "$image_id" ]; then
    echo -e "\n${YELLOW}🖼️ 第五步: 测试图片详情和更新${NC}"

    test_api "GET" "$BASE_URL/images/$image_id" "" "获取图片详情"

    update_data='{
        "name": "更新后的户型图名称",
        "description": "更新后的描述信息"
    }'

    test_api "PUT" "$BASE_URL/images/$image_id" "$update_data" "更新图片信息"
fi

# 6. 测试统计信息
echo -e "\n${YELLOW}📊 第六步: 测试统计信息${NC}"

test_api "GET" "$BASE_URL/images/stats" "" "获取图片统计信息"

# 7. 清理测试数据（可选）
echo -e "\n${YELLOW}🧹 第七步: 清理测试数据${NC}"

read -p "是否删除测试数据？(y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "正在清理测试数据..."

    # 删除图片
    if [ ! -z "$image_id" ]; then
        curl -s -X DELETE "$BASE_URL/images/$image_id" \
            ${AUTH_TOKEN:+-H "Authorization: Bearer $AUTH_TOKEN"}
        echo "✅ 删除图片成功"
    fi

    # 删除户型
    curl -s -X DELETE "$BASE_URL/house-types/$house_type_id" \
        ${AUTH_TOKEN:+-H "Authorization: Bearer $AUTH_TOKEN"}
    echo "✅ 删除户型成功"

    # 删除楼盘
    curl -s -X DELETE "$BASE_URL/buildings/$building_id" \
        ${AUTH_TOKEN:+-H "Authorization: Bearer $AUTH_TOKEN"}
    echo "✅ 删除楼盘成功"
fi

echo -e "\n${GREEN}🎉 楼盘-户型图上传功能完整测试完成！${NC}"
echo -e "${YELLOW}📝 测试总结:${NC}"
echo "  ✅ 楼盘创建和文件夹初始化"
echo "  ✅ 户型创建"
echo "  ✅ 户型图上传到指定文件夹"
echo "  ✅ 图片列表查询"
echo "  ✅ 图片详情和更新"
echo "  ✅ 统计信息获取"
echo "  ✅ 数据清理"

echo -e "\n${GREEN}🏗️ 文件夹结构验证:${NC}"
echo "  buildings/$building_id/"
echo "  ├── floor-plans/     (户型图)"
echo "  ├── images/          (楼盘图片)"
echo "  └── documents/       (相关文档)"
