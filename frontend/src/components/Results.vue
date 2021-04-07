<template>
  <div>
    <div class="row mt-2 mb-2">
      <div class="col-lg-6">
        <div class="row">
          <div class="form-inline col-6">
            <div class="form-group w-100">
              <select v-model="target" class="form-control form-control-sm">
                <option :value="t" v-for="t in targets" :key="t">
                  {{ t }}
                </option>
              </select>
              <button
                @click="getIfaces()"
                class="btn btn-sm btn-outline-primary ml-1"
                title="Update IP list"
              >
                <fa icon="sync" />
              </button>
            </div>
          </div>
          <div class="col-6 form-inline">
              <div class="ml-auto">
                <select class="form-control form-control-sm" v-model="autoScan">
                  <option value="0">Auto scan disabled</option>
                  <option value="10">Auto scan 10 secs</option>
                  <option value="20">Auto scan 20 secs</option>
                  <option value="30">Auto scan 30 secs</option>
                  <option value="60">Auto scan 1 minute</option>
                  <option value="120">Auto scan 2 inutes</option>
                </select>
                <button
                  class="btn btn-outline-primary btn-sm ml-1"
                  @click="scan()"
                  :disabled="scanning"
                >
                  SCAN
                </button>
            </div>
          </div>
        </div>
      </div>
      <div class="col-lg-6">
        <div class="row ">
          <div class="col">
            <span v-if="scanning" class="form-control-plaintext scanning">
              <fa icon="spinner" spin />
              Scanning...
            </span>
          </div>
          <div class="col form-inline ml-auto">
          <div class="form-inline ml-auto">
            <p class="form-control-plaintext">
              {{ ignore.length }} hidden
            </p>
            <button
              v-if="ignore.length > 0"
              @click="ignore = []"
              class="btn btn-primary btn-sm ml-2"
            >
              Show all
            </button>
          </div>
          </div>
        </div>
      </div>
    </div>
    <table class="table table-sm table-striped">
      <colgroup>
        <col style="text-align: centre; width: 3em" />
        <col />
        <col />
        <col />
        <col />
        <col />
        <col style="text-align: right; width: 3em" />
      </colgroup>
      <thead>
        <tr>
          <th></th>
          <th colspan="2">
            <sort-title field="Changed" v-model="sort" class="text-nowrap"
              >Changed</sort-title
            >
          </th>
          <th>
            <sort-title field="LastSeen" v-model="sort" class="text-nowrap"
              >Last Seen</sort-title
            >
          </th>
          <th>
            <sort-title field="IP" v-model="sort" class="text-nowrap"
              >IP</sort-title
            >
          </th>
          <th>
            <sort-title field="Name" v-model="sort" class="text-nowrap"
              >Name</sort-title
            >
          </th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in sortedResults" :key="item.IP">
          <th>
            <fa
              class="heart"
              icon="heart"
              :class="{ active: item.Active, inactive: !item.Active }"
            />
          </th>
          <th>{{ item.Changed | date }}</th>
          <th>{{ item.Changed | age }}</th>
          <td>{{ item.LastSeen| date }}</td>
          <td>{{ item.IP }}</td>
          <td>
            <a :href="'http://' + item.Name" target="_blank">{{ item.Name }}</a>
          </td>
          <td>
            <button
              class="btn btn-outline-secondary btn-sm"
              @click="handleIgnore(item.IP)"
              title="hide"
            >
              <fa icon="eye-slash" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import SortTitle from "./SortTitle.vue";

export default {
  components: {
    SortTitle,
  },
  filters: {
    age(v) {
      const now = new Date()
      const age = Math.round((now - v) / 1000)
      if (age > 3600) {
        return Math.floor(age/3600) + " hr"
      }
      if (age > 60) {
        return Math.floor(age/60) + " min"
      }
      return age + " sec"
    },
    date(v) {
      const now = new Date()
      if (now.getDate() == v.getDate() &&
        now.getMonth() == v.getMonth() &&
        now.getFullYear() == v.getFullYear()) {
          return v.toLocaleTimeString()
      }
      return v.toLocaleString()
    }
  },
  data() {
    return {
      targets: [],
      target: "",
      items: [],
      ignore: [],
      autoScan: 10,
      scanning: false,
      sort: {
        field: "Changed",
        asc: true,
      },
    };
  },
  computed: {
    sortedResults() {
      return this.items
        .slice()
        .filter((item) => this.ignore.indexOf(item.IP) === -1)
        .sort((a, b) => {
          if (this.sort.asc) {
            const t = a;
            a = b;
            b = t;
          }
          switch (this.sort.field) {
            case "LastSeen":
            case "Changed":
              return a[this.sort.field] - b[this.sort.field];
            case "Name":
            case "IP":
              return a[this.sort.field].localeCompare(b[this.sort.field]);
          }
        });
    },
  },
  created() {
    this.getIfaces();
  },
  watch: {
    target() {
      if (this.target && !this.scanning) {
        this.scan();
      }
    },
    autoScan() {
      if (this.autoScan > 0 && !this.scanning) {
        this.scan();
      }
    },
  },
  methods: {
    handleIgnore(ip) {
      if (this.ignore.indexOf(ip) === -1) {
        this.ignore.push(ip);
      }
    },
    handleThClick(col) {
      if (this.sortBy == col) {
        this.sortAsc = !this.sortAsc;
      } else {
        this.sortBy = col;
      }
    },
    getIfaces() {
      window.backend.GetIPv4NonLocalInterfaces().then((results, err) => {
        if (err) {
          alert(err);
          return;
        }
        this.targets = results;
        if (!this.target && this.targets.length > 0) {
          this.target = this.targets[0];
        }
      });
    },
    scan() {
      if (this.scanning) return
      this.scanning = true;
      // const now = Math.round(new Date().getTime() / 1000);
      window.backend.scan(this.target).then((results, err) => {
        this.scanning = false;
        if (this.autoScan > 0) {
          setTimeout(this.scan, 1000 * this.autoScan);
        }
        if (err) {
          alert(err);
          return;
        }
        this.items = results.map((item) => {
          item.Changed = new Date(item.Changed);
          item.LastSeen = new Date(item.LastSeen);
          return item;
        });
      });
    },
  },
};
</script>

<style scoped>
.heart {
  font-size: 1.6em;
  display: block;
  margin: auto;
}
.inactive {
  color: silver;
}
.active {
  color: green;
}
.scanning {
  color: green;
  font-weight: bold;
}
.hidden {
  float: right;
}
</style>