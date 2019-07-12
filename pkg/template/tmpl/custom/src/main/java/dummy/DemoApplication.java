package {{.PackageName}};

{{$hasDekorate := .HasDekorate}}

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;{{if $hasDekorate}}
import io.dekorate.kubernetes.annotation.KubernetesApplication;
import io.dekorate.openshift.annotation.EnableS2iBuild;
import io.dekorate.openshift.annotation.OpenshiftApplication;
{{end}}

@SpringBootApplication{{if $hasDekorate}}
@KubernetesApplication
@EnableS2iBuild
@OpenshiftApplication{{end}}
public class DemoApplication {

    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }
}
