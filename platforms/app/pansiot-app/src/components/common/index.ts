/**
 * 通用组件入口文件
 */

import CustomNavBar from './CustomNavBar/index.vue';
import PageContainer from './PageContainer/index.vue';
import Loading from './Loading/index.vue';
import EmptyState from './EmptyState/index.vue';
import NetworkError from './NetworkError/index.vue';
import PullRefresh from './PullRefresh/index.vue';
import LoadMore from './LoadMore/index.vue';

export { CustomNavBar, PageContainer, Loading, EmptyState, NetworkError, PullRefresh, LoadMore };

export default {
  CustomNavBar,
  PageContainer,
  Loading,
  EmptyState,
  NetworkError,
  PullRefresh,
  LoadMore,
};
