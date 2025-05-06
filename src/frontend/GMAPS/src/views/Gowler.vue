<template>
    <div>
        <h1>Gowler</h1>

        <input type="text" placeholder="Website..." v-model="textWebSite" />
        <button @click="crawl">
            <span v-if="isLoading">
                <i class="fa-solid fa-spinner fa-spin"></i>
            </span>
            <span v-else>
                <i class="fa-solid fa-arrow-right"></i>
            </span>
            Crawl
        </button>
        <button @click="downloadXlsx" :disabled="!gowlerData">Download XLSX</button>

        <h2>Site Information</h2>
        <div>
            <p><strong>Site:</strong> {{ gowlerData ? gowlerData.site : '' }}</p>
            <p><strong>Domain:</strong> {{ gowlerData ? gowlerData.domain : '' }}</p>
        </div>

        <h3>Site URLs</h3>
        <div class="card">
            <ul v-if="gowlerData">
                <li v-for="url in gowlerData.siteUrls" :key="url">
                    <a :href="processedSiteUrls(url)" target="_blank" rel="noopener noreferrer">{{ url }}</a>
                </li>
            </ul>
        </div>

        <h3>Other URLs</h3>
        <div class="card">
            <ul v-if="gowlerData">
                <li v-for="url in gowlerData.otherUrls" :key="url">
                    <a :href="processedSiteUrls(url)" target="_blank" rel="noopener noreferrer">{{ url }}</a>
                </li>
            </ul>
        </div>

        <h3>Telephones</h3>
        <div class="card">
            <ul v-if="gowlerData">
                <li v-for="phone in gowlerData.telephones" :key="phone">{{ phone }}</li>
            </ul>
        </div>
        <h3>Emails</h3>
        <div class="card">
            <ul v-if="gowlerData">
                <li v-for="email in gowlerData.emails" :key="email">{{ email }}</li>
            </ul>
        </div>
    </div>
</template>

<script>
    import { ref } from 'vue';
    import { useRoute } from 'vue-router';
    
    import * as XLSX from 'xlsx';
    
    export default {
        name: 'Gowler',
        components: { },
        setup() {
            const route = useRoute();
            const textWebSite = ref('');
            const textArea = ref('');
            const gowlerData = ref(null);
            const isLoading = ref(false);

            if (route.params.website) {
                textWebSite.value = route.params.website;
            }

            // Function to crawl the website and get the data
            const crawl = async () => {
                isLoading.value = true;
                
                try {
                    new URL(textWebSite.value);
                } catch (e) {
                    alert('Please enter a valid URL.');
                    isLoading.value = false;
                    return;
                }
                
                const params = new URLSearchParams(
                    {
                        website: textWebSite.value,
                    });
                const response = await fetch(`/api/gowler?${params}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
                const data = await response.json();
                gowlerData.value = data;
                
                isLoading.value = false;
            };

            // Function to process the URLs
            const processedSiteUrls = (url) => {
                let site = url;
                if (site.startsWith('/') || site.startsWith('#')) {
                    site = url.includes(gowlerData.value.site) ? url : `${gowlerData.value.domain}${url}`;
                } 
                if (site.startsWith('https://')) {
                    return site;
                } else if (site.startsWith('http://')) {
                    return site;
                } else {
                    return `https://${site}`;
                }
            };

            // Function to download CSV file
            const downloadXlsx = () => {
                const now = Date.now();
                const sheetData1 = [
                    ['Site', 'Domain'],
                    [gowlerData.value?.site || '', gowlerData.value?.domain || '']
                ];
                const sheetData2 = [
                    ['Site URLs'],
                    ...(gowlerData.value?.siteUrls ?? []).map(url => [url])
                ];
                const sheetData3 = [
                    ['Other URLs'],
                    ...(gowlerData.value?.otherUrls ?? []).map(url => [url])
                ];
                const sheetData4 = [
                    ['Telephones'],
                    ...(gowlerData.value?.telephones ?? []).map(phone => [phone])
                ];
                const sheetData5 = [
                    ['Emails'],
                    ...(gowlerData.value?.emails ?? []).map(email => [email])
                ];

                const ws1 = XLSX.utils.aoa_to_sheet(sheetData1);
                const ws2 = XLSX.utils.aoa_to_sheet(sheetData2);
                const ws3 = XLSX.utils.aoa_to_sheet(sheetData3);
                const ws4 = XLSX.utils.aoa_to_sheet(sheetData4);
                const ws5 = XLSX.utils.aoa_to_sheet(sheetData5);

                const wb = XLSX.utils.book_new();

                XLSX.utils.book_append_sheet(wb, ws1, 'Site Info');
                XLSX.utils.book_append_sheet(wb, ws2, 'Site URLs');
                XLSX.utils.book_append_sheet(wb, ws3, 'Other URLs');
                XLSX.utils.book_append_sheet(wb, ws4, 'Telephones');
                XLSX.utils.book_append_sheet(wb, ws5, 'Emails');

                const fileName = `${gowlerData.value.domain.toLowerCase()}_${now}.xlsx`;
                XLSX.writeFile(wb, fileName);
            };

            return {
                textWebSite,
                textArea,
                gowlerData,
                isLoading,

                crawl,
                processedSiteUrls,
                downloadXlsx,
            };
        },
    }
</script>

<style scoped>
    h2 {
        margin-bottom: 10px;
    }
    h3 {
        margin-top: 10px;
        margin-bottom: 5px;
    }
    ul {
        list-style-type: none;
        padding-left: 0;
    }
    li {
        margin-bottom: 5px;
    }
    .card {
        border: 1px solid #ccc;
        padding: 10px;
        margin-bottom: 20px;
        border-radius: 5px;
        background-color: #f9f9f9;
        height: 200px;
        overflow-y: auto;
    }
    .card ul {
        margin: 0;
        padding: 0;
    }
    .card li {
        margin: 0;
        padding: 0;
    }
    .card a {
        text-decoration: none;
        color: #007bff;
    }
    .card a:hover {
        text-decoration: underline;
    }
    input[type="text"] {
        width: 300px;
        padding: 5px;
        margin-bottom: 10px;
    }
</style>