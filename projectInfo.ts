import YAML from 'yaml';
import { toCanvas } from 'html-to-image';
import { message } from '@hsmos/message';
import { defineStore } from 'pinia';
import {
  getJobTree,
  addJobApi,
  deleteJobApi,
  saveJobApi,
  getJobInfoApi,
  updateJobApi,
  startJobApi,
  stopJobApi,
  debuggerJobApi,
  setFlowImg,
  checkPortApi,
  getPreviewData,
  getIntervalApi,
} from '@/api/job-flow/common';
import { API_CODE, JOB_STATUS, PERIOD_OPT, PROCESS_COMP, JOB_TYPE_OPT } from '@/constant/index';
import { useNodeStatusStore } from './nodeStatus';
import { useAppConfigStore } from './appConfig';
import bus from '@/utils/bus';
import locales from '@/locales';
import type {
  JobForm,
  TreeRoot,
  ActiveJob,
  NodeItem,
  LineItem,
  BenthosYamlJson,
  TreeLeafStatusItem,
  projectInfoItem,
} from '@/views/job-flow/types';
import type { Message } from 'element-plus';
import type { Ref } from 'vue';

const { t } = locales.global;
// 调度单位对应简称
const periodCValueMap = Object.values(PERIOD_OPT).reduce((map: CommonObject<string>, el) => {
  const { value, cvalue } = el as { value: string; cvalue: string };
  map[value] = cvalue as string;
  return map;
}, {});
interface State {
  jobHasChange: boolean;
  selectNode: Partial<JobForm>;
  rootNodeId: string;
  projectInfo: Partial<projectInfoItem>;
  projectTreeData: TreeRoot[];
  activeJob: Partial<JobForm>;
  jobData: ActiveJob;
  jobDataSets: CommonObject<any>;
  debuggerIng: boolean;
  isDebugger: boolean;
  mainLinks: string[];
  envVarList: CommonObject<string>[];
  debugResult: any;
  nodeYamlMap: CommonObject<CommonObject<any>>;
  benthosYamlJson: BenthosYamlJson;
  flowConentDom: null | HTMLElement;
}
type receiveInfoFunc = (msg: string, isSuccess: boolean, isOver: boolean) => void;
type activeConnectFunc = (node: NodeItem) => void;
const nodeStatus = useNodeStatusStore();
const { inlayEnv, mode, getProjectAttr, userInfo } = useAppConfigStore();
const namespace = getProjectAttr();
const IS_ENG = mode !== 'DEV' && inlayEnv() === 'eng';
const IS_PORTAL = mode !== 'DEV' && inlayEnv() === 'portal';
const $message = message(bus) as Message;
export const useProjectInfoStore = defineStore('projectInfo', {
  state: (): State => ({
    jobHasChange: false, // 作业信息是否发生改变
    selectNode: {}, // 工程树选中的节点
    rootNodeId: '', //工程树根节点id
    projectInfo: {}, //项目信息,用于删除作业时传给eng关闭tab页
    debugResult: null, // 调试结果，用于存放generator
    debuggerIng: false, // 是否正在调试
    isDebugger: false, // 是否开启调试
    mainLinks: [], // 主链路数据
    envVarList: [], // 环境变量列表
    projectTreeData: [
      {
        id: '0',
        name: '信息系统集成',
        icon: 'info-system',
        level: 1,
        children: [],
      },
    ],
    jobDataSets: {},
    jobData: {},
    activeJob: {}, // 树节点选中【作业】后画布展示的作业
    nodeYamlMap: {
      input: {},
      transform: {},
      output: {},
    },
    benthosYamlJson: { input: {}, pipeline: {}, output: {}, logger: {} },
    flowConentDom: null, // 画布内容区域的dom元素
  }),
  actions: {
    // 作业数据集合
    updateJobDataSets(key: string, data: any) {
      this.jobDataSets[key] = data;
    },
    removeJobDataSets(key: string) {
      this.jobDataSets[key] = null;
    },
    getJobDataSets(key: string) {
      return this.jobDataSets[key];
    },
    getTreeRootNodeId() {
      const src = window.location.hash;
      if (src.includes('?')) {
        const arr = src.split('?')[1].split('&');
        arr.forEach((item) => {
          const m = item.split('=');
          if (m[0] == 'mmenu_id') {
            this.rootNodeId = m[1];
          }
        });
      }
    },
    /**
     * 根据状态获取类型和文本
     * @param status 状态值
     * @returns
     */
    getNodeStatus(status: number): TreeLeafStatusItem {
      const obj = JOB_STATUS[status] || {};
      const statusType = obj.statusType || '';
      const statusText = obj.label || '';
      const statusTextEN = obj.statusTextEN || '';
      const statusColor = obj.statusColor || '';
      return { statusType, statusText, statusTextEN, statusColor };
    },
    async getJobList(jobClass = '', operatorObj?: { name: string; type: 'delete' | 'add' }) {
      const params = {
        namespace,
        job_class: jobClass,
      };
      const { data } = await getJobTree(params);
      if (!data?.length) {
        // $message.info('工程树暂无服务节点！');
        return;
      }
      data.forEach((item: any) => {
        item.canAdd = true;
        item.icon = 'tree-server-node';
        item.level = 2;
        if (!item.children) {
          item.children = [];
        } else {
          item.children.forEach((el: JobForm) => {
            el.icon = 'tree-job';
            el.node = item.id;
            el.id = el.code;
            el.isLeaf = true;

            // 根据【状态】设置状态【文本】和【类型】
            if (el.status !== null && el.status !== undefined) {
              el = Object.assign(el, this.getNodeStatus(el.status));
            }
            // 更新选中节点的数据
            if (this.selectNode.isLeaf && this.selectNode.id === el.code) {
              this.selectNode = el;
            }
          });
        }
      });
      this.projectTreeData[0].children = data;
      if (IS_ENG) {
        !this.selectNode.id && (this.selectNode = data[0]);
        this.statusLaneInfo(operatorObj);
      }
    },
    async getJobInfo(code: string) {
      const params = {
        code,
        namespace,
      };
      const { code: rescode, data } = await getJobInfoApi(params);
      if (rescode === API_CODE.success && data && data.code) {
        const { canvas, content, ...jobInfo } = data;
        this.activeJob = jobInfo;
        this.jobData = {
          canvas,
          content,
        };
        if (data?.content.benthos_yaml) {
          this.benthosYamlJson = Object.assign(
            { input: {}, pipeline: {}, output: {}, logger: {} },
            YAML.parse(data?.content.benthos_yaml)
          );
        } else {
          this.benthosYamlJson = { input: {}, pipeline: {}, output: {}, logger: {} };
        }
        if (data?.content.node_yaml_map) {
          this.nodeYamlMap = Object.assign({}, data?.content.node_yaml_map);
        } else {
          this.nodeYamlMap = {
            input: {},
            transform: {},
            output: {},
          };
        }
        if (data?.content.node_data_map) {
          this.jobDataSets = Object.assign({}, data?.content.node_data_map);
        } else {
          this.jobDataSets = {};
        }
        if (data?.content?.canvas_param) {
          this.envVarList = [...data.content.canvas_param];
        } else {
          this.envVarList = [];
        }
      } else {
        $message.error(t('requesterror.' + rescode));
        this.jobData = {};
        this.activeJob = {};
        this.selectNode = {};
      }
    },
    // 新增作业
    async add(result: JobForm) {
      // let params = {
      //   describe: result.describe,
      //   name: result.name,
      //   s_node: result.node,
      //   schedule: result.schedule,
      //   schedule_cycle: String(result.schedule_cycle),
      //   schedule_unit: result.schedule_unit,
      //   template: result.template,
      //   job_type: result.job_type || 'timing',
      //   job_class: result.job_class,
      //   namespace,
      //   create_user: 'nickName' in userInfo ? (userInfo.nickName as string) : '',
      // };
      const params = Object.assign(result, {
        s_node: result.node,
        schedule_cycle: String(result.schedule_cycle),
        job_type: result.job_type || 'timing',
        namespace,
        create_user: 'nickName' in userInfo ? (userInfo.nickName as string) : '',
      });
      const { code, data, message } = await addJobApi(params);
      if (code === API_CODE.success && data) {
        $message.success(message);
        // 新增成功后并选中节点
        this.selectNode = { ...result, code: data.id, id: data.id, isLeaf: true, status: 0 };
        await this.getJobList(result.job_class, { name: result.name, type: 'add' });
        if (IS_ENG) {
          bus.post({
            type: 'project-tree-refresh',
            data: {
              id: this.rootNodeId,
              activeId: data.id,
              parentId: this.projectTreeData[0].children[0].id,
            },
          });
          bus.emit('project-tree-add'); //发消息通知预览界面更新
        }
        this.getJobInfo(data.id);
      }
      if (!data && message) {
        $message.error(message);
      }
    },
    // 删除作业
    async remove(jobItem?: JobForm) {
      if (jobItem) {
        const params = {
          code: jobItem.code!,
          namespace,
        };
        const { code, data, message: msg } = await deleteJobApi(params);
        if (data) {
          $message.success(msg);
          if (IS_ENG || IS_PORTAL) {
            bus.post({
              type: 'project-tree-refresh',
              data: { id: this.rootNodeId },
            });
            bus.emit('project-tree-remove-by-pre', jobItem.code!);
          }
          await this.getJobList(jobItem.job_class, { name: jobItem.name!, type: 'delete' });
        } else {
          $message.error(t('requesterror.' + code));
        }
      } else if (this.activeJob.code) {
        const params = {
          code: this.activeJob.code,
          namespace,
        };
        const { code, data, message } = await deleteJobApi(params);
        if (code === API_CODE.success && data) {
          $message.success(message);
          if (IS_ENG || IS_PORTAL) {
            await this.getJobList(this.activeJob.job_class, {
              name: this.activeJob.name!,
              type: 'delete',
            });
            bus.post({
              type: 'project-tree-refresh',
              data: { id: this.rootNodeId },
            });
            // 通知eng框架关闭信息系统tab页
            bus.post({
              type: 'close-tab-data',
              closeType: 'specify-tab',
              data: { allowDelete: true, ...this.projectInfo },
            });
            bus.emit('ioit-preview-refresh');
          }
          // 取消选中节点
          this.selectNode = {};
        } else if (code !== API_CODE.success) {
          $message.error(t('requesterror.' + code));
        } else {
          return false;
        }
      } else {
        return false;
      }
    },
    // 更新作业【有编辑作业弹窗时调用】
    async update(result: JobForm, nodeList: NodeItem[], lineList: LineItem[]) {
      this.jobHasChange = false;
      const yamlItem = await this.generatorYaml(lineList, nodeList);
      // const params = {
      //   code: result.code || '',
      //   describe: result.describe,
      //   name: result.name,
      //   s_node: result.node,
      //   schedule: result.schedule,
      //   schedule_cycle: String(result.schedule_cycle),
      //   schedule_unit: result.schedule_unit,
      //   schedule_details: result.schedule_details,
      //   job_class: result.job_class,
      //   is_http: this.isAutomaticTrigeer(result.job_type || 'timing'),
      //   job_type: result.job_type || 'timing',
      //   template: result.template,
      //   canvas: {
      //     nodeList: nodeList.map((node) => ({
      //       ...node,
      //       ...nodeStatus.getNodeStatusMap(node.id),
      //     })),
      //     lineList,
      //   },
      //   content: {
      //     benthos_yaml: yamlItem,
      //   },
      //   namespace,
      //   create_user: 'nickName' in userInfo ? (userInfo.nickName as string) : '',
      // };
      const params = Object.assign(result, {
        code: result.code || '',
        s_node: result.node,
        schedule_cycle: String(result.schedule_cycle),
        is_http: this.isAutomaticTrigeer(result.job_type || 'timing'),
        job_type: result.job_type || 'timing',
        canvas: {
          nodeList: nodeList.map((node) => ({
            ...node,
            ...nodeStatus.getNodeStatusMap(node.id),
          })),
          lineList,
        },
        content: {
          benthos_yaml: yamlItem,
        },
        namespace,
        create_user: 'nickName' in userInfo ? (userInfo.nickName as string) : '',
      });
      const { code, data, message } = await updateJobApi(params);
      if (code === API_CODE.success && data) {
        $message.success(message);
        this.selectNode = Object.assign(this.selectNode, result, { icon: 'tree-job' });
        // if (IS_ENG) {
        //   bus.emit(
        //     'project-tree-update',
        //     JSON.stringify({
        //       code: result.code,
        //       selectNode: this.selectNode,
        //     })
        //   );
        // }
        if (IS_ENG || IS_PORTAL) {
          bus.post({
            type: 'project-tree-refresh',
            data: { id: this.rootNodeId },
          });
          bus.emit('project-tree-job-save', result.code);
        }
        await this.getJobList(result.job_class);
        // 更新成功后并选中节点,获取其数据
        result.code && (await this.getJobInfo(result.code));
      } else if (code !== API_CODE.success) {
        $message.error(t('requesterror.' + code));
      }
    },
    // 保存作业
    async save(nodeList: NodeItem[], lineList: LineItem[], form: any, loading: Ref<boolean>) {
      if (this.activeJob.code) {
        loading.value = true;
        // this.setFlowToImg(this.activeJob.code);
        this.jobHasChange = false;
        // 需要更新节点的状态信息
        this.jobData.canvas = {
          nodeList: nodeList.map((node) => ({
            ...node,
            ...nodeStatus.getNodeStatusMap(node.id),
          })),
          lineList,
        };
        if (!this.jobData?.content) {
          this.jobData.content = {};
        }
        this.jobData.content.node_data_map = this.jobDataSets;
        this.jobData.content!.node_yaml_map = this.nodeYamlMap;
        this.jobData.content!.canvas_param = this.envVarList;
        this.jobData.content!.benthos_yaml = await this.generatorYaml(lineList, nodeList);
        this.activeJob = Object.assign(this.activeJob, form, {
          job_type: form.job_type || 'timing',
        });
        const params = {
          ...this.activeJob,
          ...this.jobData,
          is_http: this.isAutomaticTrigeer(form.job_type || 'timing'),
          schedule_cycle: String(this.activeJob.schedule_cycle),
          create_user: 'nickName' in userInfo ? (userInfo.nickName as string) : '',
        } as any;
        try {
          const { code, data, message } = await saveJobApi(params);
          if (code === API_CODE.success && data) {
            $message.success(message);
            this.selectNode = Object.assign(this.selectNode, form);
            if (this.activeJob.code) {
              await this.getJobInfo(this.activeJob.code);
            }
            if (IS_ENG || IS_PORTAL) {
              bus.emit('project-tree-job-save', this.activeJob.code);
            }
          } else if (code !== API_CODE.success) {
            $message.error(t('requesterror.' + code));
          }
        } catch {
          loading.value = false;
        }
        loading.value = false;
      }
    },
    // 组合流程yaml获取数据
    previewData(data: any, processYml: CommonObject<any>) {
      const yaml = {
        input: {
          generate: {
            mapping: `root = ${JSON.stringify(data)}`,
            interval: '0s',
            count: 1,
          },
        },
        pipeline: {
          processors: Array.isArray(processYml) ? processYml : [processYml],
        },
        output: {
          broker: {
            outputs: [
              {
                file: {
                  codec: 'all-bytes',
                  path: '/usr/local/hsm-os/data/hsm-io-it/data/benthos/data/dataPreview.json',
                },
              },
            ],
          },
        },
      };
      return getPreviewData(YAML.stringify(yaml));
    },
    // 将画布转为base64
    setFlowToImg(id: string) {
      const dom = this.flowConentDom;
      if (!dom) return '';
      const canvasInfo = {
        scale: window.devicePixelRatio,
        height: dom.clientHeight,
        width: dom.clientWidth,
        backgroundColor: undefined,
      };
      toCanvas(dom, canvasInfo).then((canvas) => {
        const url = canvas.toDataURL();
        setFlowImg({ id, url });
      });
    },
    // 设置容器元素
    setContentElement(dom: HTMLElement | null) {
      this.flowConentDom = dom;
    },
    // 设置节点的yaml映射关系
    setNodeYaml(node: NodeItem, yamlItem: any) {
      this.nodeYamlMap[node.flowNodeType][node.id] = yamlItem;
    },
    // 获取节点的yaml映射
    getNodeYaml(node: NodeItem) {
      return this.nodeYamlMap[node.flowNodeType][node.id] || '';
    },
    // 删除节点的yaml映射关系
    deleteNodeYaml(node: NodeItem) {
      delete this.nodeYamlMap[node.flowNodeType][node.id];
    },
    // 通过id获取节点的类型
    getNodeTypes(id: string) {
      return Object.keys(this.nodeYamlMap).find((key) => this.nodeYamlMap[key][id]);
    },
    // 是否为自动触发，自动触发通过benthos调度，否则通过后端cron调度
    // isAutomaticTrigeer(nodeList: NodeItem[]) {
    //   return (
    //     (this.hasTargetCompInput(nodeList, [
    //       'HTTPClient',
    //       'ERPClient',
    //       'HTTPServer',
    //       'WEBService',
    //     ]) &&
    //       this.hasScheduleConf()) ||
    //     this.hasTargetCompInput(nodeList, ['MySQL', 'DataBase'])
    //   );
    // },
    isAutomaticTrigeer(jobType: JobForm['job_type']) {
      if (jobType === 'realtime') {
        return true;
      } else {
        return IS_ENG;
      }
    },
    // 获取完整链路信息
    getIntegrityLink(
      nodeMap: CommonObject<NodeItem>,
      nodeLinksMap: CommonObject<CommonObject<LineItem[]>>
    ) {
      // 获取完成链路信息，会将同一输入节点的分支作为一条链路输出
      // const activeInputItem = {
      //   branch: '',
      //   index: 0,
      // };
      // const branchMap: CommonObject<string[][]> = {};
      // const getNextLink = (input: string, links: string[]): string[][] => {
      //   if (input) {
      //     links.push(input);
      //     const inputNextNodes = (nodeLinksMap[input] && nodeLinksMap[input].outputLink) || [];
      //     if (!inputNextNodes.length) return branchMap[activeInputItem.branch];
      //     if (inputNextNodes.length > 1) {
      //       const branchLinks = inputNextNodes.map((_) => [...links]);
      //       branchMap[activeInputItem.branch].splice(activeInputItem.index, 1, ...branchLinks);
      //       inputNextNodes.forEach((line) => {
      //         getNextLink(line.target, branchMap[activeInputItem.branch][activeInputItem.index]);
      //         activeInputItem.index++;
      //       });
      //     } else {
      //       getNextLink(
      //         inputNextNodes[0].target,
      //         branchMap[activeInputItem.branch]
      //           ? branchMap[activeInputItem.branch][activeInputItem.index]
      //           : links
      //       );
      //     }
      //   }
      //   return branchMap[activeInputItem.branch];
      // };
      // 同一输入的整条链路元素，采用深度优先遍历，节点顺序并不是按照链路走
      const getNextLink = (input: string, links: string[], isBranch = false): string[] => {
        if (input) {
          links.push(input);
          const inputNextNodes = (nodeLinksMap[input] && nodeLinksMap[input].outputLink) || [];
          if (!inputNextNodes.length && nodeMap[input]?.flowNodeType !== 'output') {
            // 删除无输出端的分支路线
            if (isBranch) {
              links.pop();
              while (
                nodeMap[links[links.length - 1]].flowNodeType === 'transform' &&
                nodeLinksMap[links[links.length - 1]].outputLink.length <= 1
              ) {
                links.pop();
              }
              // 如果主分支没有输出端直接将链路置为空数组
            } else {
              links.length = 0;
            }
          } else {
            inputNextNodes.forEach((line) => {
              getNextLink(line.target.cell, links, inputNextNodes.length > 1);
            });
          }
        }
        return links;
      };
      // 获取到所有链路后过滤掉没有输出端的链路
      const linkNodeIds: string[][] = [];
      Object.keys(this.nodeYamlMap.input).forEach((input) => {
        const nextLinks = getNextLink(input, []);
        if (!nextLinks.length) return false;
        if (!linkNodeIds.length) linkNodeIds.push(nextLinks);
        else {
          let hasMerge = false;
          linkNodeIds.forEach((links, index) => {
            const mergeLinks = [...new Set([...links, ...nextLinks])];
            // 通过set去重数组
            if (mergeLinks.length !== links.length + nextLinks.length) {
              const sortLinks = [];
              for (let i = 0, j = 0, len = links.length, len2 = nextLinks.length; i < len; i++) {
                const item = links[i];
                const mergeIdx = nextLinks.indexOf(item);
                if (mergeIdx !== -1) {
                  sortLinks.push(...nextLinks.slice(j, mergeIdx));
                  j = mergeIdx + 1;
                }
                sortLinks.push(item);
                // 最后判断是否遍历完毕，确保其他分支的数据被添加
                if (i + 1 === len && j !== len2) {
                  sortLinks.push(...nextLinks.slice(j));
                }
              }
              linkNodeIds[index] = sortLinks;
              hasMerge = true;
            }
          });
          if (!hasMerge) {
            linkNodeIds.push(nextLinks);
          }
        }
      });
      // 将多输入链路进行合并，只有当两条链路存在重合节点才被识别为同一条链路，最终获取出主链路
      const mainLineNodes = (linkNodeIds || []).reduce((arr, next) => {
        if (next.length > arr.length) {
          arr = next;
        }
        return arr;
      }, []);
      this.mainLinks = mainLineNodes;
      return mainLineNodes;
    },
    // 生成yaml文件
    async generatorYaml(lineList: LineItem[], nodeList: NodeItem[]) {
      const { schExp } = await this.getInterval();
      const processComp = nodeList.find((el) => PROCESS_COMP.includes(el.value) && el.isConfig);
      // 处理一个组件可以作为输入输出的情况
      if (processComp) {
        this.benthosYamlJson['input'] = this.getNodeYaml({
          ...processComp,
          flowNodeType: 'input',
        });
        this.benthosYamlJson['output'] = this.getNodeYaml({
          ...processComp,
          flowNodeType: 'output',
        });
      } else {
        const nodeMap: CommonObject<NodeItem> = {}; // 节点map
        const nodeLinksMap = nodeList.reduce((linkMap, node) => {
          const { id } = node;
          nodeMap[id] = node;
          const inputLink = lineList.filter((link) => link.target.cell === id);
          const outputLink = lineList.filter((link) => link.source.cell === id);
          linkMap[id] = { inputLink, outputLink };
          return linkMap;
        }, {} as CommonObject<CommonObject<LineItem[]>>); // 节点连接map，分in和out
        this.benthosYamlJson = { input: {}, pipeline: {}, output: {}, logger: {} };
        const mainLineNodes = this.getIntegrityLink(nodeMap, nodeLinksMap);
        let outputFieldConvert = null;
        // 是否存在转换组件
        const hasMySqlInput = this.hasTargetCompInput(nodeList, ['MySQL', 'DataBase']);
        // const hasBatchInput = this.hasTargetCompInput(nodeList, ['HTTPServer'], 'input');
        const hasMySqlOutput = this.hasTargetCompInput(nodeList, ['DataBase'], 'output');
        let occupiedOutput = false; // 是否占用输出组件
        let lastProcessors: CommonObject<any> | null = null; // 最后一个process
        mainLineNodes.reduce((collection, next, idx) => {
          const nodeInfo = nodeMap[next];
          const nodeType = nodeInfo.flowNodeType;
          const targetYaml = this.nodeYamlMap[nodeType][next];
          if (!targetYaml) return collection;
          // 是否存在最后一个流程组件配置，存在就在更新完毕后添加至process中
          if (targetYaml['lastProcessors']) {
            lastProcessors = targetYaml['lastProcessors'];
          }
          if (nodeType === 'input') {
            // 存在mysql的话输入转为定时器，input改到process中
            // 统一数据格式
            const targetYamlArr = Array.isArray(targetYaml) ? targetYaml : [targetYaml];
            if (hasMySqlInput) {
              // 存在mysql输出的话进行输出展示
              if (hasMySqlOutput) {
                collection['input'] = {
                  label: 'mysql',
                  ...targetYamlArr[0]['branch']['processors'][0],
                };
              } else {
                collection['input'] = this.generatorIntervalYaml(schExp);
                const nowProcessors = collection['pipeline']['processors'] || [];
                // 数据源会存在分支数据
                const hasBranch = nowProcessors.findLastIndex(
                  (process: CommonObject<any>) => process['branch']
                );
                if (hasBranch !== -1) {
                  nowProcessors.splice(hasBranch + 1, 0, ...targetYamlArr);
                } else {
                  collection['input'] = this.generatorIntervalYaml(schExp);
                  const nowProcessors = collection['pipeline']['processors'] || [];
                  // 数据源会存在分支数据
                  const hasBranch = nowProcessors.findLastIndex(
                    (process: CommonObject<any>) => process['branch']
                  );
                  if (hasBranch !== -1) {
                    nowProcessors.splice(hasBranch + 1, 0, ...targetYamlArr);
                  } else {
                    nowProcessors.push(...targetYamlArr);
                  }
                  collection['pipeline'] = {
                    processors: nowProcessors,
                  };
                  // nowProcessors.push(...targetYamlArr);
                }
                // collection['pipeline'] = {
                //   processors: nowProcessors,
                // };
              }
            } else {
              // 由于object中key会乱序，为了在输入后可以方便处理，这里添加hasInputAndProcess，为true得话说明input为输入数据，process为处理数据，分别添加至不同得yaml分类
              if (targetYamlArr[0]['hasInputAndProcess'] || targetYamlArr[0]['hasInputAndOutput']) {
                collection['input'] = targetYamlArr[0]['input'];
                if (targetYamlArr[0]['processors']) {
                  const nowProcessors = collection['pipeline']['processors'] || [];
                  collection['pipeline'] = {
                    processors: nowProcessors.concat(
                      Array.isArray(targetYamlArr[0]['processors'])
                        ? targetYamlArr[0]['processors']
                        : [targetYamlArr[0]['processors']]
                    ),
                  };
                }
              } else {
                collection['input'] = targetYamlArr[0];
              }
              if (targetYamlArr[0]['hasInputAndOutput']) {
                collection['output'] = targetYamlArr[0]['output'];
                occupiedOutput = true;
              }
            }
          } else if (nodeType === 'transform') {
            const nowProcessors = collection['pipeline']['processors'] || [];
            collection['pipeline'] = {
              processors: nowProcessors.concat(
                Array.isArray(targetYaml) ? targetYaml : [targetYaml]
              ),
            };
          } else {
            // 占用输出得情况
            if (occupiedOutput) {
              let outputYaml = targetYaml;
              // 存在粘连process的时候添加process然后进行输出
              if ('synechiaProcess' in outputYaml) {
                this.addYamlToArr(
                  this.benthosYamlJson['pipeline'],
                  'processors',
                  outputYaml['synechiaProcess']
                );
                outputYaml = outputYaml['output'];
              }
              // 存在分支得情况需继续判断，http客户端情况
              if ('broker' in outputYaml) {
                const {
                  broker: { outputs },
                } = outputYaml;
                if (outputs[1] && 'http_client' in outputs[1]) {
                  outputs[0] = {
                    label: 'http',
                    http: outputs[1]['http_client'],
                  };
                }
                outputYaml = outputs[0];
              }

              this.addYamlToArr(this.benthosYamlJson['pipeline'], 'processors', outputYaml);
              if (idx === mainLineNodes.length - 1 && lastProcessors) {
                this.addYamlToArr(this.benthosYamlJson['pipeline'], 'processors', lastProcessors);
              }
              return collection;
            }
            // 是否存在输出字段，存在则需要进行转换
            if (targetYaml.outputField) {
              const copyYaml = { ...targetYaml };
              outputFieldConvert = copyYaml.outputField;
              delete copyYaml.outputField;
              collection[nodeType] = copyYaml;
            } else {
              if ('synechiaProcess' in targetYaml) {
                collection[nodeType] = targetYaml['output'];
                this.addYamlToArr(
                  this.benthosYamlJson['pipeline'],
                  'processors',
                  targetYaml['synechiaProcess']
                );
              } else {
                collection[nodeType] = targetYaml;
              }

              // if (Array.isArray(targetYaml)) {
              //   const nowProcessors = collection['pipeline']['processors'] || [];
              //   collection['pipeline'] = {
              //     processors: nowProcessors.concat([targetYaml[0]]),
              //   };
              //   collection[nodeType] = targetYaml[1];
              // } else {
              //   collection[nodeType] = targetYaml;
              // }
            }
          }
          if (idx === mainLineNodes.length - 1 && lastProcessors) {
            this.addYamlToArr(this.benthosYamlJson['pipeline'], 'processors', lastProcessors);
          }
          return collection;
        }, this.benthosYamlJson);
        // 如果存在说明需要进行转换，添加转换yaml配置
        if (outputFieldConvert) {
          this.addYamlToArr(this.benthosYamlJson['pipeline'], 'processors', outputFieldConvert);
        }
        // if (hasBatchInput) {
        //   this.addYamlToArr(this.benthosYamlJson['pipeline'], 'processors', {
        //     mapping: 'root = {"code":"1","msg":"ok","data":"true"}',
        //   });
        // }
        // 如果当前作业是周期调度，且为http客户端则需要生成调度资源的yaml
        // if (
        //   this.hasScheduleConf() &&
        //   this.hasTargetCompInput(nodeList, ['HTTPClient', 'ERPClient'])
        // ) {
        //   this.benthosYamlJson = Object.assign(
        //     {},
        //     this.benthosYamlJson,
        //     this.geteratorJobIntervalYaml()
        //   );
        // }
      }
      // 生成日志的yaml
      this.benthosYamlJson.logger = {
        level: 'ALL',
        format: 'json',
        add_timestamp: true,
        file: {
          path: `/usr/local/hsm-os/data/hsm-io-it/data/benthos/json/${this.activeJob.code}.json`,
          rotate: true,
        },
      };
      return YAML.stringify(this.benthosYamlJson);
    },
    // 初始化或追加yaml在数组中
    addYamlToArr(data: CommonObject<any[]>, key: string, val: any) {
      if (!Array.isArray(val)) val = [val];
      data[key] ? data[key].push(...val) : (data[key] = val);
    },
    // 生成作业定时任务的yaml
    // geteratorJobIntervalYaml() {
    //   if (!this.hasScheduleConf()) {
    //     return null;
    //   }
    //   const { schedule_cycle, schedule_unit } = this.activeJob;
    //   const item = {
    //     rate_limit_resources: [
    //       {
    //         label: 'foo',
    //         local: {
    //           count: 1,
    //           interval: schedule_cycle + periodCValueMap[schedule_unit as string],
    //         },
    //       },
    //     ],
    //   };
    //   return item;
    // },
    getCount() {
      if (this.activeJob.schedule === JOB_TYPE_OPT.once.value) {
        return 1;
      }
      if (
        this.activeJob.schedule === JOB_TYPE_OPT.period.value &&
        this.activeJob?.period_detail?.endType === 'times'
      ) {
        return this.activeJob?.period_detail.times;
      }
      return;
    },
    async getInterval() {
      const params = {
        code: this.activeJob.code,
        namespace,
      };
      const { code, data } = await getIntervalApi(params);
      if (code === API_CODE.success && data) {
        return data;
      } else if (code !== API_CODE.success) {
        $message.error(t('requesterror.' + code));
        return data;
      }
    },
    // 生成作业定时器yaml，用于控制mysql等args_mapping: hasWhere ? '[0]' : '[1]',
    // dsn: `sqlserver://${username}:${password}@${host}:${port}?database=${selectedDataBaseName.value}`,
    //是否以数字结尾
    // export const isNumberEnding = {
    //   test: (val: string) => /(\d+)$/.test(val),
    //   exec: (val: string) => /(\d+)$/.exec(val),
    //   replace: (val: string, rVal = '') => val.replace(/(\d+)$/, rVal),
    // };
    generatorIntervalYaml(schExp: string) {
      console.log(
        'generatorIntervalYaml====00',
        /^(\d+)/.test(schExp),
        `${JSON.stringify(schExp)}`
      );
      const item = {
        label: 'interval',
        generate: {
          // interval: /^(\d+)/.test(schExp) ? `\'${schExp}\'` : `${schExp}`,
          interval: '5555 55',
          // interval: schExp,
          count: this.getCount(),
          mapping: 'root = {}',
        },
      };
      return item;
    },
    // 包含周期调度配置
    hasScheduleConf() {
      const { schedule, schedule_cycle, schedule_unit } = this.activeJob;
      if (!schedule || !schedule_cycle || !schedule_unit) {
        return false;
      }
      if (schedule !== '2') {
        return false;
      }
      return true;
    },
    // 是否包含目标组件
    hasTargetCompInput(
      nodeList: NodeItem[],
      comp: string | string[],
      flowNodeType = 'input',
      type?: string
    ): boolean | NodeItem {
      const comps = typeof comp === 'string' ? [comp] : comp;
      const targetComp = nodeList.find((node) => {
        return (
          node.flowNodeType === flowNodeType &&
          comps.includes(node.value) &&
          this.mainLinks.includes(node.id)
        );
      });
      if (type === 'nodeItem' && targetComp) {
        return targetComp;
      } else {
        return !!targetComp;
      }
    },

    /**
     * 启动作业
     */
    async start() {
      this.debugResult = null;
      this.isDebugger = false;
      this.debuggerIng = false;
      if (this.activeJob.code) {
        const params = {
          code: this.activeJob.code,
          namespace,
        };
        // HTTPServer时判断端口是否占用
        const nodeMap: CommonObject<NodeItem> = {}; // 节点map
        const nodeLinksMap =
          this.jobData.canvas?.nodeList.reduce((linkMap, node) => {
            const { id } = node;
            nodeMap[id] = node;
            const inputLink =
              this.jobData.canvas?.lineList.filter((link) => link.target.cell === id) || [];
            const outputLink =
              this.jobData.canvas?.lineList.filter((link) => link.source.cell === id) || [];
            linkMap[id] = { inputLink, outputLink };
            return linkMap;
          }, {} as CommonObject<CommonObject<LineItem[]>>) || {}; // 节点连接map，分in和out
        this.getIntegrityLink(nodeMap, nodeLinksMap);
        const targetComp = this.hasTargetCompInput(
          this.jobData.canvas?.nodeList || [],
          ['Api'],
          'input',
          'nodeItem'
        );
        if (typeof targetComp !== 'boolean' && 'port' in (targetComp?.attrs || {})) {
          const { data: portData, code } = await checkPortApi({
            batchport: targetComp.attrs.port.toString(),
          });
          if (!portData) {
            $message.error(t('requesterror.' + code));
            return;
          }
        }
        const { code, data, message } = await startJobApi(params);
        if (code === API_CODE.success && data) {
          if (IS_ENG || IS_PORTAL) {
            bus.emit(
              'project-tree-status',
              JSON.stringify({
                code: this.activeJob.code,
                type: {
                  ...this.getNodeStatus(1),
                  status: 1,
                },
              })
            );
            bus.post({
              type: 'project-tree-refresh',
              data: {
                id: this.rootNodeId,
                activeId: this.activeJob.code,
                parentId: this.projectTreeData[0].children[0].id,
              },
            });
          }
          await this.getJobList(this.activeJob.job_class || '');
          return { data, message };
        } else if (code !== API_CODE.success) {
          $message.error(t('requesterror.' + code));
        }
      }
    },

    /**
     * 停止作业
     */
    async stop() {
      if (this.activeJob.code) {
        const params = {
          code: this.activeJob.code,
          namespace,
        };
        const { code, data, message } = await stopJobApi(params);
        if (code === API_CODE.success && data) {
          if (IS_ENG || IS_PORTAL) {
            bus.emit(
              'project-tree-status',
              JSON.stringify({
                code: this.activeJob.code,
                type: {
                  ...this.getNodeStatus(0),
                  status: 0,
                },
              })
            );
            bus.post({
              type: 'project-tree-refresh',
              data: {
                id: this.rootNodeId,
                activeId: this.activeJob.code,
                parentId: this.projectTreeData[0].children[0].id,
              },
            });
          }
          await this.getJobList();
          return { data, message };
        } else if (code !== API_CODE.success) {
          $message.error(t('requesterror.' + code));
        }
      }
    },
    // 作业是否发生变化
    modifyJobStatus(hasChange: boolean) {
      if (hasChange !== this.jobHasChange) {
        this.jobHasChange = hasChange;
      }
    },
    // 调式作业
    debuggerJob(nodeList: NodeItem[], infoFunc: receiveInfoFunc, activeConnect: activeConnectFunc) {
      this.isDebugger = true;
      // 调试作业
      if (!this.debugResult) {
        this.debugResult = this.singleNodeDebugger(nodeList, infoFunc, activeConnect);
      }
    },
    // 停止调试
    stopDebugger() {
      this.debugResult = null;
      this.isDebugger = false;
      this.debuggerIng = false;
      $message.info('调试结束');
    },
    nextDebugger() {
      if (this.debuggerIng) {
        $message.warning('正在调试中！请调试完毕后继续下一步');
        return;
      }
      if (this.debugResult) {
        this.debugResult.next();
      }
    },
    // 利用generator进行各个节点的调试
    singleNodeDebugger: function* (
      nodeList: NodeItem[],
      infoFunc: receiveInfoFunc,
      activeConnect: activeConnectFunc
    ) {
      for (let i = 0, len = nodeList.length; i < len; i++) {
        this.debuggerIng = true;
        const node = nodeList[i];
        activeConnect(node);
        nodeStatus.editNodeStatusMap(node.id, 'active', '调试中');
        const nodeYaml = this.getNodeYaml(node);
        const { outputdata, inputdata } = node;
        const params = {
          datas: { outputdata, inputdata },
          yaml: nodeYaml,
        };
        yield debuggerJobApi(params).then((res) => {
          const { data, msg } = res;
          nodeStatus.editNodeStatusMap(node.id, data ? 'success' : 'error', msg);
          const debugInfo = `节点：${node.label} 节点数据：${JSON.stringify(
            params.datas.outputdata
          )} 调试信息:${msg}`;
          infoFunc(debugInfo, data, i === len - 1);
          if (i === len - 1) this.debugResult.next();
          this.debuggerIng = false;
        });
      }
      this.stopDebugger();
    },
    // ENG状态栏信息展示
    statusLaneInfo(operatorObj?: { name: string; type: 'delete' | 'add' }) {
      let statusInfoMap = {
        all: 0,
        start: 0,
      };
      this.projectTreeData[0].children.forEach((el) => {
        statusInfoMap = el.children.reduce((map, next) => {
          if (next.status) {
            map['start'] += 1;
          }
          map['all'] += 1;
          return map;
        }, statusInfoMap);
      });
      const operatorInfos = operatorObj
        ? [
            {
              type: 'success',
              content: `${operatorObj.type === 'delete' ? '删除' : '新增'}作业<${
                operatorObj.name
              }>成功`,
            },
          ]
        : [];
      const dataItem: CommonObject<{ type: string; content: string }[]> = {
        number: Object.keys(statusInfoMap).map((el) => ({
          type: 'info',
          content: `信息系统集成作业${el === 'start' ? '运行成功' : ''}总数: ${
            statusInfoMap[el as keyof typeof statusInfoMap]
          }`,
        })),
      };
      if (operatorInfos.length) {
        dataItem['operation'] = operatorInfos;
      }
      bus.post({
        type: 'eng-state-show',
        data: dataItem,
      });
    },
  },
});
