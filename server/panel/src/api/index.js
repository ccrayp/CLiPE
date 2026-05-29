import { createCrudApi } from './crudApi'

export const apiMap = {
  users: createCrudApi({
    basePath: '/users',
    listKey: 'users',
  }),
  hosts: createCrudApi({
    basePath: '/hosts',
    listKey: 'hosts',
  }),
  services: createCrudApi({
    basePath: '/services',
    listKey: 'services',
  }),
  rules: createCrudApi({
    basePath: '/rules',
    listKey: 'rules',
  }),
  policies: createCrudApi({
    basePath: '/policies',
    listKey: 'policies',
  }),
  policyContents: createCrudApi({
    basePath: '/policy-contents',
    listKey: 'policy_contents',
    buildUpdatePath: ({ policy_id, service_id }) =>
      `/policy-contents/${policy_id}/${service_id}`,
    buildDeletePath: ({ policy_id, service_id }) =>
      `/policy-contents/${policy_id}/${service_id}`,
  }),
  requests: createCrudApi({
    basePath: '/requests',
    listKey: 'requests',
  }),
  decisions: createCrudApi({
    basePath: '/decisions',
    listKey: 'decisions',
  }),
}
