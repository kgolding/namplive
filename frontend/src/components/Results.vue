<template>
  <div>
    <div class="row  mt-2 mb-2">
      <div class="col">
        <span v-if="scanning">
          <fa icon="spinner" spin/>
          Scanning...
        </span>
      </div>
      <div class="col hidden text-right">
        {{ ignore.length }} hidden
        <button
          v-if="ignore.length > 0"
          @click="ignore = []"
          class="btn btn-primary btn-sm"
        >
          Show all
        </button>
      </div>
    </div>
    <table class="table table-sm table-striped">
      <colgroup>
        <col style="text-align: centre; width: 3em">
        <col>
        <col>
        <col>
        <col>
        <col style="text-align: right; width: 3em">
      </colgroup>
      <thead>
        <tr>
          <th></th>
          <th>
            <sort-title field="Changed" v-model="sort" class="text-nowrap"
              >Changed</sort-title
            >
          </th>
          <th>
            <sort-title field="LastNotSeen" v-model="sort" class="text-nowrap"
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
          <th>{{ item.Changed.toLocaleString() }}</th>
          <td>{{ item.LastSeen.toLocaleString() }}</td>
          <!-- <td>{{ item.LastSeen.toLocaleString() }}</td> -->
          <td>{{ item.IP }}</td>
          <td>
            <a :href="'http://' + item.Name" target="_blank">{{ item.Name }}</a>
          </td>
          <td>
            <button
              class="btn btn-outline-primary btn-sm"
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
  data() {
    return {
      items: [],
      ignore: [],
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
    this.scan();
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
    scan() {
      this.scanning = true;
      // const now = Math.round(new Date().getTime() / 1000);
      window.backend.scan().then((results, err) => {
        this.scanning = false;
        setTimeout(this.scan, 10000);
        if (err) {
          alert(err);
          return;
        }
        this.items = results.map((item) => {
          item.Changed = new Date(item.Changed);
          item.LastSeen = new Date(item.LastSeen);
          // item.Age = now - item.LastSeen.getTime();
          // item.Age = now - Math.round(item.LastSeen.getTime() / 1000);
          // item.Active = item.LastSeen > item.LastNotSeen
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
.hidden {
  float: right;
}
</style>