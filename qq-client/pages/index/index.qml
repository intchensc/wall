<scroll-view class="container">
  <button qq:if="{{!hasUserInfo && canIUse}}" open-type="getUserInfo" bindgetuserinfo="getUserInfo">获取头像昵称 </button>

  <block qq:if="{{hasUserInfo}}">
    <view class="form">
      <!-- <text>{{app.globalData.userInfo}}</text> -->
      <form bindsubmit="touGao">
        <textarea placeholder="你需要为你投稿带来的消极影响负责:)" name="textarea"/>
        
        <checkbox-group bindchange="checkboxChange">
          <label class="checkbox" qq:for="{{items}}">
            <checkbox value="{{item.name}}" checked="{{item.checked}}" />
            {{item.value}}
          </label>
        </checkbox-group>
        <button form-type="submit">提交</button>
        <text onclick="admin">成为管理员</text>
      </form>
    </view>
  </block>
</scroll-view>

